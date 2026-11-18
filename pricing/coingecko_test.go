package pricing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoinGeckoPrecos(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/simple/price" {
			t.Errorf("path = %q, esperado /simple/price", r.URL.Path)
		}
		if got := r.URL.Query().Get("symbols"); got != "btc,eth,xyz" {
			t.Errorf("symbols = %q, esperado btc,eth,xyz", got)
		}
		if got := r.URL.Query().Get("vs_currencies"); got != "brl" {
			t.Errorf("vs_currencies = %q, esperado brl", got)
		}
		if got := r.Header.Get("x-cg-demo-api-key"); got != "chave-teste" {
			t.Errorf("header da api key = %q, esperado chave-teste", got)
		}
		_, _ = w.Write([]byte(`{"btc":{"brl":368000.5},"eth":{"brl":10300.123456789}}`))
	}))
	defer srv.Close()

	c := NovoCoinGecko(srv.URL, "chave-teste", "BRL")
	precos, err := c.Precos(context.Background(), []string{"BTC", "ETH", "XYZ"})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if c.Moeda() != "brl" {
		t.Errorf("Moeda() = %q, esperado brl", c.Moeda())
	}
	if len(precos) != 2 {
		t.Fatalf("len(precos) = %d, esperado 2: %v", len(precos), precos)
	}
	if precos["BTC"].String() != "368000.5" {
		t.Errorf("BTC = %s, esperado 368000.5", precos["BTC"])
	}
	if precos["ETH"].String() != "10300.123456789" {
		t.Errorf("ETH = %s, esperado 10300.123456789", precos["ETH"])
	}
}

func TestCoinGeckoSemSimbolosNaoFazRequisicao(t *testing.T) {
	c := NovoCoinGecko("http://127.0.0.1:1", "", "usd")
	precos, err := c.Precos(context.Background(), nil)
	if err != nil || len(precos) != 0 {
		t.Fatalf("esperava mapa vazio sem erro, obtido %v, %v", precos, err)
	}
}

func TestCoinGeckoStatusDeErro(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c := NovoCoinGecko(srv.URL, "", "usd")
	if _, err := c.Precos(context.Background(), []string{"BTC"}); err == nil {
		t.Fatal("esperava erro para status 429")
	}
}
