package models

import (
	"errors"
	"regexp"
	"strings"
)

// CryptoAsset representa um ativo da carteira.
type CryptoAsset struct {
	ID       int     `json:"id"`
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

var simboloValido = regexp.MustCompile(`^[A-Z0-9]{1,10}$`)

// Normalizar padroniza os campos do ativo (ex.: " btc " -> "BTC").
func (a *CryptoAsset) Normalizar() {
	a.Symbol = strings.ToUpper(strings.TrimSpace(a.Symbol))
}

// Validar verifica se os campos do ativo são válidos. Chame Normalizar antes.
func (a CryptoAsset) Validar() error {
	if !simboloValido.MatchString(a.Symbol) {
		return errors.New("symbol deve ter de 1 a 10 caracteres alfanumericos")
	}
	if a.Quantity <= 0 {
		return errors.New("quantity deve ser maior que zero")
	}
	if a.Price < 0 {
		return errors.New("price nao pode ser negativo")
	}
	return nil
}
