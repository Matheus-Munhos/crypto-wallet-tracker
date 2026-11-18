package controllers

import (
	"crypto-tracker/models"
	"crypto-tracker/repository"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// AssetHandler agrupa os handlers HTTP da carteira.
type AssetHandler struct {
	repo repository.AssetRepository
}

func NovoAssetHandler(repo repository.AssetRepository) *AssetHandler {
	return &AssetHandler{repo: repo}
}

// CadastrarAtivo trata POST /wallet.
func (h *AssetHandler) CadastrarAtivo(w http.ResponseWriter, r *http.Request) {
	var asset models.CryptoAsset
	if err := lerJSON(w, r, &asset); err != nil {
		escreverErro(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	asset.Normalizar()
	if err := asset.Validar(); err != nil {
		escreverErro(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Criar(r.Context(), &asset); err != nil {
		log.Printf("erro ao salvar ativo: %v", err)
		escreverErro(w, http.StatusInternalServerError, "Erro ao salvar no banco de dados")
		return
	}

	escreverJSON(w, http.StatusCreated, asset)
}

// ListarAtivos trata GET /wallet.
func (h *AssetHandler) ListarAtivos(w http.ResponseWriter, r *http.Request) {
	assets, err := h.repo.Listar(r.Context())
	if err != nil {
		log.Printf("erro ao listar ativos: %v", err)
		escreverErro(w, http.StatusInternalServerError, "Erro ao listar ativos")
		return
	}
	escreverJSON(w, http.StatusOK, assets)
}

// BuscarAtivo trata GET /wallet/{id}.
func (h *AssetHandler) BuscarAtivo(w http.ResponseWriter, r *http.Request) {
	id, ok := lerID(w, r)
	if !ok {
		return
	}

	asset, err := h.repo.Buscar(r.Context(), id)
	if err != nil {
		h.tratarErroRepositorio(w, err, "Erro ao buscar ativo")
		return
	}
	escreverJSON(w, http.StatusOK, asset)
}

// AtualizarAtivo trata PUT /wallet/{id}.
func (h *AssetHandler) AtualizarAtivo(w http.ResponseWriter, r *http.Request) {
	id, ok := lerID(w, r)
	if !ok {
		return
	}

	var asset models.CryptoAsset
	if err := lerJSON(w, r, &asset); err != nil {
		escreverErro(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	asset.ID = id
	asset.Normalizar()
	if err := asset.Validar(); err != nil {
		escreverErro(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Atualizar(r.Context(), asset); err != nil {
		h.tratarErroRepositorio(w, err, "Erro ao atualizar ativo")
		return
	}
	escreverJSON(w, http.StatusOK, asset)
}

// RemoverAtivo trata DELETE /wallet/{id}.
func (h *AssetHandler) RemoverAtivo(w http.ResponseWriter, r *http.Request) {
	id, ok := lerID(w, r)
	if !ok {
		return
	}

	if err := h.repo.Remover(r.Context(), id); err != nil {
		h.tratarErroRepositorio(w, err, "Erro ao remover ativo")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// tratarErroRepositorio responde 404 para ativo inexistente e 500 para os demais erros.
func (h *AssetHandler) tratarErroRepositorio(w http.ResponseWriter, err error, mensagem string) {
	if errors.Is(err, repository.ErrNaoEncontrado) {
		escreverErro(w, http.StatusNotFound, "Ativo nao encontrado")
		return
	}
	log.Printf("%s: %v", mensagem, err)
	escreverErro(w, http.StatusInternalServerError, mensagem)
}

// lerID extrai e valida o parâmetro {id} da URL.
func lerID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		escreverErro(w, http.StatusBadRequest, "id invalido")
		return 0, false
	}
	return id, true
}

// CalcularTotal trata GET /wallet/total.
func (h *AssetHandler) CalcularTotal(w http.ResponseWriter, r *http.Request) {
	total, err := h.repo.TotalInvestido(r.Context())
	if err != nil {
		log.Printf("erro ao calcular total: %v", err)
		escreverErro(w, http.StatusInternalServerError, "Erro ao calcular total")
		return
	}
	escreverJSON(w, http.StatusOK, map[string]float64{"total_invested": total})
}
