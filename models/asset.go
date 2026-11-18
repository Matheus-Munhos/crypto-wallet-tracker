package models

import (
	"errors"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

// CryptoAsset representa um ativo da carteira.
//
// Quantity e Price usam decimal para não perder precisão. No JSON eles saem
// como string (ex.: "0.5") e na entrada aceitam número ou string.
type CryptoAsset struct {
	ID       int             `json:"id"`
	Symbol   string          `json:"symbol"`
	Quantity decimal.Decimal `json:"quantity"`
	Price    decimal.Decimal `json:"price"`
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
	if !a.Quantity.IsPositive() {
		return errors.New("quantity deve ser maior que zero")
	}
	if a.Price.IsNegative() {
		return errors.New("price nao pode ser negativo")
	}
	return nil
}
