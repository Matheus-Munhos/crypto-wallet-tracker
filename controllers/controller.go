package controller

import (
	"context"
	"crypto-tracker/database"
	"encoding/json"
	"net/http"
)

type CryptoAsset struct {
	ID       int     `json:"id"`
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

func CadastrarAtivo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var asset CryptoAsset
	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		http.Error(w, `{"error": "JSON invalido"}`, http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO crypto_assets (symbol, quantity, price) 
		VALUES ($1, $2, $3) 
		RETURNING id;
	`

	err = database.DB.QueryRow(context.Background(), query, asset.Symbol, asset.Quantity, asset.Price).Scan(&asset.ID)
	if err != nil {
		http.Error(w, `{"error": "Erro ao salvar no banco de dados"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(asset)

}

func ListarArquivos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, _ := database.DB.Query(context.Background(), "SELECT id, symbol, quantity, price FROM crypto_assets")
	defer rows.Close()

	var assets []CryptoAsset

	for rows.Next() {
		var asset CryptoAsset
		rows.Scan(&asset.ID, &asset.Symbol, &asset.Quantity, &asset.Price)
		assets = append(assets, asset)

	}
	json.NewEncoder(w).Encode(assets)
}

func CalcularTotal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `SELECT COALESCE(SUM(quantity * price), 0) FROM crypto_assets;`

	var totalInvestido float64

	err := database.DB.QueryRow(context.Background(), query).Scan(&totalInvestido)
	if err != nil {
		http.Error(w, `{"error": "Erro ao calcular total"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"Total_investido": totalInvestido})

}
