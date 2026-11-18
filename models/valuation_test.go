package models

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCalcularAvaliacao(t *testing.T) {
	d := decimal.RequireFromString

	assets := []CryptoAsset{
		{ID: 1, Symbol: "BTC", Quantity: d("0.5"), Price: d("100")}, // investido 50, vale 60
		{ID: 2, Symbol: "ETH", Quantity: d("2"), Price: d("25")},    // investido 50, vale 40
		{ID: 3, Symbol: "XYZ", Quantity: d("10"), Price: d("1")},    // sem cotação
		{ID: 4, Symbol: "XYZ", Quantity: d("5"), Price: d("1")},     // sem cotação (repetido)
	}
	precos := map[string]decimal.Decimal{"BTC": d("120"), "ETH": d("20")}

	av := CalcularAvaliacao(assets, precos, "usd")

	verificar := func(nome string, obtido decimal.Decimal, esperado string) {
		t.Helper()
		if !obtido.Equal(d(esperado)) {
			t.Errorf("%s = %s, esperado %s", nome, obtido, esperado)
		}
	}

	if av.Currency != "usd" {
		t.Errorf("Currency = %q, esperado usd", av.Currency)
	}
	verificar("TotalInvested", av.TotalInvested, "100")
	verificar("CurrentValue", av.CurrentValue, "100")
	verificar("ProfitLoss", av.ProfitLoss, "0")
	verificar("ProfitLossPct", av.ProfitLossPct, "0")

	if len(av.Assets) != 2 {
		t.Fatalf("len(Assets) = %d, esperado 2", len(av.Assets))
	}
	verificar("BTC ProfitLoss", av.Assets[0].ProfitLoss, "10")
	verificar("BTC ProfitLossPct", av.Assets[0].ProfitLossPct, "20")
	verificar("ETH ProfitLoss", av.Assets[1].ProfitLoss, "-10")
	verificar("ETH ProfitLossPct", av.Assets[1].ProfitLossPct, "-20")

	if !reflect.DeepEqual(av.MissingPrices, []string{"XYZ"}) {
		t.Errorf("MissingPrices = %v, esperado [XYZ]", av.MissingPrices)
	}
}

func TestCalcularAvaliacaoCarteiraVazia(t *testing.T) {
	av := CalcularAvaliacao(nil, nil, "brl")

	if !av.TotalInvested.IsZero() || !av.ProfitLossPct.IsZero() {
		t.Errorf("esperava totais zerados, obtido %+v", av)
	}
	if av.Assets == nil || av.MissingPrices == nil {
		t.Error("Assets e MissingPrices devem ser listas vazias, nao nil (viram [] no JSON)")
	}
}
