package controller

import (
	"context"
	"crypto-tracker/database"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type CryptoAsset struct {
	ID       int     `json:"id"`
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

var simboloValido = regexp.MustCompile(`^[A-Z0-9]{1,10}$`)

// validar normaliza o símbolo e verifica se os campos do ativo são válidos.
func (a *CryptoAsset) validar() string {
	a.Symbol = strings.ToUpper(strings.TrimSpace(a.Symbol))
	if !simboloValido.MatchString(a.Symbol) {
		return "symbol deve ter de 1 a 10 caracteres alfanumericos"
	}
	if a.Quantity <= 0 {
		return "quantity deve ser maior que zero"
	}
	if a.Price < 0 {
		return "price nao pode ser negativo"
	}
	return ""
}

func escreverJSON(w http.ResponseWriter, status int, dados any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(dados); err != nil {
		log.Printf("erro ao escrever resposta JSON: %v", err)
	}
}

func escreverErro(w http.ResponseWriter, status int, mensagem string) {
	escreverJSON(w, status, map[string]string{"error": mensagem})
}

func CadastrarAtivo(w http.ResponseWriter, r *http.Request) {
	var asset CryptoAsset
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
		escreverErro(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if msg := asset.validar(); msg != "" {
		escreverErro(w, http.StatusBadRequest, msg)
		return
	}

	query := `
		INSERT INTO crypto_assets (symbol, quantity, price)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	err := database.DB.QueryRow(context.Background(), query, asset.Symbol, asset.Quantity, asset.Price).Scan(&asset.ID)
	if err != nil {
		log.Printf("erro ao salvar ativo: %v", err)
		escreverErro(w, http.StatusInternalServerError, "Erro ao salvar no banco de dados")
		return
	}

	escreverJSON(w, http.StatusCreated, asset)
}

func ListarArquivos(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(context.Background(), "SELECT id, symbol, quantity, price FROM crypto_assets ORDER BY id")
	if err != nil {
		log.Printf("erro ao listar ativos: %v", err)
		escreverErro(w, http.StatusInternalServerError, "Erro ao listar ativos")
		return
	}
	defer rows.Close()

	assets := []CryptoAsset{}
	for rows.Next() {
		var asset CryptoAsset
		if err := rows.Scan(&asset.ID, &asset.Symbol, &asset.Quantity, &asset.Price); err != nil {
			log.Printf("erro ao ler ativo: %v", err)
			escreverErro(w, http.StatusInternalServerError, "Erro ao listar ativos")
			return
		}
		assets = append(assets, asset)
	}
	if err := rows.Err(); err != nil {
		log.Printf("erro ao percorrer ativos: %v", err)
		escreverErro(w, http.StatusInternalServerError, "Erro ao listar ativos")
		return
	}

	escreverJSON(w, http.StatusOK, assets)
}

func CalcularTotal(w http.ResponseWriter, r *http.Request) {
	query := `SELECT COALESCE(SUM(quantity * price), 0) FROM crypto_assets;`

	var totalInvestido float64
	if err := database.DB.QueryRow(context.Background(), query).Scan(&totalInvestido); err != nil {
		log.Printf("erro ao calcular total: %v", err)
		escreverErro(w, http.StatusInternalServerError, "Erro ao calcular total")
		return
	}

	escreverJSON(w, http.StatusOK, map[string]float64{"Total_investido": totalInvestido})
}
