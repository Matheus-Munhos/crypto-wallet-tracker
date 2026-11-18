package models

import (
	"sort"

	"github.com/shopspring/decimal"
)

// AtivoAvaliado é um ativo com sua cotação atual e o resultado (lucro/prejuízo).
type AtivoAvaliado struct {
	CryptoAsset
	CurrentPrice  decimal.Decimal `json:"current_price"`
	Invested      decimal.Decimal `json:"invested"`
	CurrentValue  decimal.Decimal `json:"current_value"`
	ProfitLoss    decimal.Decimal `json:"profit_loss"`
	ProfitLossPct decimal.Decimal `json:"profit_loss_percent"`
}

// Avaliacao é o resumo da carteira a preços atuais.
//
// Os totais consideram apenas os ativos com cotação disponível; os símbolos
// sem cotação aparecem em MissingPrices.
type Avaliacao struct {
	Currency      string          `json:"currency"`
	TotalInvested decimal.Decimal `json:"total_invested"`
	CurrentValue  decimal.Decimal `json:"current_value"`
	ProfitLoss    decimal.Decimal `json:"profit_loss"`
	ProfitLossPct decimal.Decimal `json:"profit_loss_percent"`
	Assets        []AtivoAvaliado `json:"assets"`
	MissingPrices []string        `json:"missing_prices"`
}

var cem = decimal.NewFromInt(100)

// CalcularAvaliacao cruza os ativos da carteira com as cotações atuais.
func CalcularAvaliacao(assets []CryptoAsset, precos map[string]decimal.Decimal, moeda string) Avaliacao {
	av := Avaliacao{
		Currency:      moeda,
		Assets:        []AtivoAvaliado{},
		MissingPrices: []string{},
	}

	faltando := map[string]bool{}
	for _, a := range assets {
		preco, ok := precos[a.Symbol]
		if !ok {
			faltando[a.Symbol] = true
			continue
		}

		investido := a.Quantity.Mul(a.Price)
		atual := a.Quantity.Mul(preco)
		resultado := atual.Sub(investido)

		av.Assets = append(av.Assets, AtivoAvaliado{
			CryptoAsset:   a,
			CurrentPrice:  preco,
			Invested:      investido,
			CurrentValue:  atual,
			ProfitLoss:    resultado,
			ProfitLossPct: percentual(resultado, investido),
		})

		av.TotalInvested = av.TotalInvested.Add(investido)
		av.CurrentValue = av.CurrentValue.Add(atual)
	}

	av.ProfitLoss = av.CurrentValue.Sub(av.TotalInvested)
	av.ProfitLossPct = percentual(av.ProfitLoss, av.TotalInvested)

	for s := range faltando {
		av.MissingPrices = append(av.MissingPrices, s)
	}
	sort.Strings(av.MissingPrices)

	return av
}

// percentual retorna parte/total * 100 com 2 casas decimais (0 se total for zero).
func percentual(parte, total decimal.Decimal) decimal.Decimal {
	if total.IsZero() {
		return decimal.Zero
	}
	return parte.Div(total).Mul(cem).Round(2)
}
