package pricing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Provider busca a cotação atual de uma lista de símbolos.
type Provider interface {
	// Precos retorna o preço atual de cada símbolo (em maiúsculas).
	// Símbolos sem cotação ficam de fora do mapa.
	Precos(ctx context.Context, simbolos []string) (map[string]decimal.Decimal, error)
	// Moeda retorna a moeda das cotações (ex.: "usd", "brl").
	Moeda() string
}

const URLPadraoCoinGecko = "https://api.coingecko.com/api/v3"

// CoinGecko implementa Provider usando a API pública da CoinGecko.
type CoinGecko struct {
	baseURL string
	apiKey  string
	moeda   string
	client  *http.Client
}

// NovoCoinGecko cria o cliente. apiKey é opcional (chave "demo" da CoinGecko).
func NovoCoinGecko(baseURL, apiKey, moeda string) *CoinGecko {
	return &CoinGecko{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		moeda:   strings.ToLower(moeda),
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *CoinGecko) Moeda() string { return c.moeda }

func (c *CoinGecko) Precos(ctx context.Context, simbolos []string) (map[string]decimal.Decimal, error) {
	precos := map[string]decimal.Decimal{}
	if len(simbolos) == 0 {
		return precos, nil
	}

	minusculos := make([]string, len(simbolos))
	for i, s := range simbolos {
		minusculos[i] = strings.ToLower(s)
	}

	params := url.Values{}
	params.Set("symbols", strings.Join(minusculos, ","))
	params.Set("vs_currencies", c.moeda)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/simple/price?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("x-cg-demo-api-key", c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar a CoinGecko: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CoinGecko respondeu com status %d", resp.StatusCode)
	}

	// Formato da resposta: {"btc": {"usd": 62755}, "eth": {"usd": 1756.95}}
	var corpo map[string]map[string]decimal.Decimal
	if err := json.NewDecoder(resp.Body).Decode(&corpo); err != nil {
		return nil, fmt.Errorf("resposta invalida da CoinGecko: %w", err)
	}

	for simbolo, cotacoes := range corpo {
		if preco, ok := cotacoes[c.moeda]; ok {
			precos[strings.ToUpper(simbolo)] = preco
		}
	}
	return precos, nil
}
