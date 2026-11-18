package models

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNormalizar(t *testing.T) {
	a := CryptoAsset{Symbol: "  btc "}
	a.Normalizar()
	if a.Symbol != "BTC" {
		t.Errorf("Symbol = %q, esperado %q", a.Symbol, "BTC")
	}
}

func TestValidar(t *testing.T) {
	d := decimal.RequireFromString

	casos := []struct {
		nome    string
		asset   CryptoAsset
		temErro bool
	}{
		{"valido", CryptoAsset{Symbol: "BTC", Quantity: d("0.5"), Price: d("100")}, false},
		{"preco zero permitido", CryptoAsset{Symbol: "ETH", Quantity: d("1"), Price: d("0")}, false},
		{"symbol vazio", CryptoAsset{Symbol: "", Quantity: d("1"), Price: d("1")}, true},
		{"symbol longo", CryptoAsset{Symbol: "ABCDEFGHIJK", Quantity: d("1"), Price: d("1")}, true},
		{"symbol com caractere invalido", CryptoAsset{Symbol: "BT-C", Quantity: d("1"), Price: d("1")}, true},
		{"quantity zero", CryptoAsset{Symbol: "BTC", Quantity: d("0"), Price: d("1")}, true},
		{"quantity negativa", CryptoAsset{Symbol: "BTC", Quantity: d("-1"), Price: d("1")}, true},
		{"price negativo", CryptoAsset{Symbol: "BTC", Quantity: d("1"), Price: d("-0.01")}, true},
	}

	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			err := c.asset.Validar()
			if (err != nil) != c.temErro {
				t.Errorf("Validar() erro = %v, esperava erro: %v", err, c.temErro)
			}
		})
	}
}
