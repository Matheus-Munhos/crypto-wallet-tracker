package controllers_test

import (
	"context"
	"crypto-tracker/controllers"
	"crypto-tracker/models"
	"crypto-tracker/repository"
	"crypto-tracker/routes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

// repoMemoria é um AssetRepository em memória para os testes.
type repoMemoria struct {
	dados  map[int]models.CryptoAsset
	proxID int
	falha  error
}

func novoRepoMemoria() *repoMemoria {
	return &repoMemoria{dados: map[int]models.CryptoAsset{}, proxID: 1}
}

func (r *repoMemoria) Criar(_ context.Context, a *models.CryptoAsset) error {
	if r.falha != nil {
		return r.falha
	}
	a.ID = r.proxID
	r.proxID++
	r.dados[a.ID] = *a
	return nil
}

func (r *repoMemoria) Listar(context.Context) ([]models.CryptoAsset, error) {
	if r.falha != nil {
		return nil, r.falha
	}
	lista := []models.CryptoAsset{}
	for _, a := range r.dados {
		lista = append(lista, a)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].ID < lista[j].ID })
	return lista, nil
}

func (r *repoMemoria) Buscar(_ context.Context, id int) (models.CryptoAsset, error) {
	a, ok := r.dados[id]
	if !ok {
		return models.CryptoAsset{}, repository.ErrNaoEncontrado
	}
	return a, nil
}

func (r *repoMemoria) Atualizar(_ context.Context, a models.CryptoAsset) error {
	if _, ok := r.dados[a.ID]; !ok {
		return repository.ErrNaoEncontrado
	}
	r.dados[a.ID] = a
	return nil
}

func (r *repoMemoria) Remover(_ context.Context, id int) error {
	if _, ok := r.dados[id]; !ok {
		return repository.ErrNaoEncontrado
	}
	delete(r.dados, id)
	return nil
}

func (r *repoMemoria) TotalInvestido(context.Context) (decimal.Decimal, error) {
	if r.falha != nil {
		return decimal.Zero, r.falha
	}
	total := decimal.Zero
	for _, a := range r.dados {
		total = total.Add(a.Quantity.Mul(a.Price))
	}
	return total, nil
}

// precosFixos é um pricing.Provider com cotações fixas.
type precosFixos struct {
	precos map[string]decimal.Decimal
	falha  error
}

func (p precosFixos) Precos(context.Context, []string) (map[string]decimal.Decimal, error) {
	return p.precos, p.falha
}

func (p precosFixos) Moeda() string { return "usd" }

func novoServidor(repo *repoMemoria, precos precosFixos) http.Handler {
	return routes.ConfigurarRotas(controllers.NovoAssetHandler(repo, precos))
}

func requisitar(t *testing.T, h http.Handler, metodo, caminho, corpo string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(metodo, caminho, strings.NewReader(corpo))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func verificarStatus(t *testing.T, rec *httptest.ResponseRecorder, esperado int) {
	t.Helper()
	if rec.Code != esperado {
		t.Fatalf("status = %d, esperado %d (corpo: %s)", rec.Code, esperado, rec.Body.String())
	}
}

func TestCadastrarAtivo(t *testing.T) {
	repo := novoRepoMemoria()
	h := novoServidor(repo, precosFixos{})

	rec := requisitar(t, h, http.MethodPost, "/wallet/", `{"symbol":" btc ","quantity":"0.5","price":100}`)
	verificarStatus(t, rec, http.StatusCreated)

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}

	var criado models.CryptoAsset
	if err := json.NewDecoder(rec.Body).Decode(&criado); err != nil {
		t.Fatal(err)
	}
	if criado.ID != 1 || criado.Symbol != "BTC" || !criado.Quantity.Equal(decimal.RequireFromString("0.5")) {
		t.Errorf("ativo criado inesperado: %+v", criado)
	}
}

func TestCadastrarAtivoInvalido(t *testing.T) {
	h := novoServidor(novoRepoMemoria(), precosFixos{})

	casos := map[string]string{
		"json malformado":    `{"symbol":`,
		"quantity negativa":  `{"symbol":"BTC","quantity":-1,"price":1}`,
		"symbol muito longo": `{"symbol":"ABCDEFGHIJK","quantity":1,"price":1}`,
		"quantity texto":     `{"symbol":"BTC","quantity":"abc","price":1}`,
	}
	for nome, corpo := range casos {
		t.Run(nome, func(t *testing.T) {
			rec := requisitar(t, h, http.MethodPost, "/wallet/", corpo)
			verificarStatus(t, rec, http.StatusBadRequest)
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("erro deve ser JSON, Content-Type = %q", ct)
			}
		})
	}
}

func TestCadastrarAtivoErroNoBanco(t *testing.T) {
	repo := novoRepoMemoria()
	repo.falha = errors.New("banco fora do ar")
	h := novoServidor(repo, precosFixos{})

	rec := requisitar(t, h, http.MethodPost, "/wallet/", `{"symbol":"BTC","quantity":1,"price":1}`)
	verificarStatus(t, rec, http.StatusInternalServerError)
}

func TestListarAtivosVazioRetornaListaVazia(t *testing.T) {
	h := novoServidor(novoRepoMemoria(), precosFixos{})

	rec := requisitar(t, h, http.MethodGet, "/wallet/", "")
	verificarStatus(t, rec, http.StatusOK)
	if corpo := strings.TrimSpace(rec.Body.String()); corpo != "[]" {
		t.Errorf("corpo = %s, esperado []", corpo)
	}
}

func TestCrudCompleto(t *testing.T) {
	h := novoServidor(novoRepoMemoria(), precosFixos{})

	verificarStatus(t, requisitar(t, h, http.MethodPost, "/wallet/", `{"symbol":"ETH","quantity":2,"price":10}`), http.StatusCreated)

	rec := requisitar(t, h, http.MethodGet, "/wallet/1", "")
	verificarStatus(t, rec, http.StatusOK)

	rec = requisitar(t, h, http.MethodPut, "/wallet/1", `{"symbol":"eth","quantity":3,"price":12}`)
	verificarStatus(t, rec, http.StatusOK)
	var atualizado models.CryptoAsset
	_ = json.NewDecoder(rec.Body).Decode(&atualizado)
	if atualizado.ID != 1 || !atualizado.Quantity.Equal(decimal.NewFromInt(3)) {
		t.Errorf("ativo atualizado inesperado: %+v", atualizado)
	}

	rec = requisitar(t, h, http.MethodGet, "/wallet/total", "")
	verificarStatus(t, rec, http.StatusOK)
	if corpo := strings.TrimSpace(rec.Body.String()); corpo != `{"total_invested":"36"}` {
		t.Errorf("total = %s", corpo)
	}

	verificarStatus(t, requisitar(t, h, http.MethodDelete, "/wallet/1", ""), http.StatusNoContent)
	verificarStatus(t, requisitar(t, h, http.MethodGet, "/wallet/1", ""), http.StatusNotFound)
}

func TestRotasComIDInvalidoOuInexistente(t *testing.T) {
	h := novoServidor(novoRepoMemoria(), precosFixos{})

	verificarStatus(t, requisitar(t, h, http.MethodGet, "/wallet/abc", ""), http.StatusBadRequest)
	verificarStatus(t, requisitar(t, h, http.MethodGet, "/wallet/0", ""), http.StatusBadRequest)
	verificarStatus(t, requisitar(t, h, http.MethodPut, "/wallet/99", `{"symbol":"BTC","quantity":1,"price":1}`), http.StatusNotFound)
	verificarStatus(t, requisitar(t, h, http.MethodDelete, "/wallet/99", ""), http.StatusNotFound)
}

func TestAvaliarCarteira(t *testing.T) {
	repo := novoRepoMemoria()
	precos := precosFixos{precos: map[string]decimal.Decimal{"BTC": decimal.NewFromInt(150)}}
	h := novoServidor(repo, precos)

	requisitar(t, h, http.MethodPost, "/wallet/", `{"symbol":"BTC","quantity":2,"price":100}`)

	rec := requisitar(t, h, http.MethodGet, "/wallet/valuation", "")
	verificarStatus(t, rec, http.StatusOK)

	var av models.Avaliacao
	if err := json.NewDecoder(rec.Body).Decode(&av); err != nil {
		t.Fatal(err)
	}
	if !av.CurrentValue.Equal(decimal.NewFromInt(300)) || !av.ProfitLossPct.Equal(decimal.NewFromInt(50)) {
		t.Errorf("avaliacao inesperada: %+v", av)
	}
}

func TestAvaliarCarteiraProvedorIndisponivel(t *testing.T) {
	h := novoServidor(novoRepoMemoria(), precosFixos{falha: errors.New("timeout")})

	rec := requisitar(t, h, http.MethodGet, "/wallet/valuation", "")
	verificarStatus(t, rec, http.StatusBadGateway)
}
