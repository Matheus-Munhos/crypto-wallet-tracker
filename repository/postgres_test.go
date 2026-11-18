package repository_test

import (
	"context"
	"crypto-tracker/database"
	"crypto-tracker/models"
	"crypto-tracker/repository"
	"errors"
	"os"
	"testing"

	"github.com/shopspring/decimal"
)

// Teste de integração: só roda quando TEST_DATABASE_URL está definida.
// ATENÇÃO: a tabela crypto_assets desse banco é apagada a cada execução.
func TestPostgresRepository(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("defina TEST_DATABASE_URL para rodar os testes de integracao")
	}

	ctx := context.Background()
	pool, err := database.Conectar(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, "TRUNCATE crypto_assets RESTART IDENTITY"); err != nil {
		t.Fatal(err)
	}

	repo := repository.NovoPostgresRepository(pool)
	d := decimal.RequireFromString

	btc := models.CryptoAsset{Symbol: "BTC", Quantity: d("0.1"), Price: d("3")}
	if err := repo.Criar(ctx, &btc); err != nil {
		t.Fatal(err)
	}
	eth := models.CryptoAsset{Symbol: "ETH", Quantity: d("0.2"), Price: d("3")}
	if err := repo.Criar(ctx, &eth); err != nil {
		t.Fatal(err)
	}

	lista, err := repo.Listar(ctx)
	if err != nil || len(lista) != 2 {
		t.Fatalf("Listar() = %v, %v", lista, err)
	}

	// 0.1*3 + 0.2*3 = 0.9 exato (com float64 daria 0.9000000000000001)
	total, err := repo.TotalInvestido(ctx)
	if err != nil || !total.Equal(d("0.9")) {
		t.Fatalf("TotalInvestido() = %s, %v; esperado 0.9", total, err)
	}

	eth.Quantity = d("1.5")
	if err := repo.Atualizar(ctx, eth); err != nil {
		t.Fatal(err)
	}
	obtido, err := repo.Buscar(ctx, eth.ID)
	if err != nil || !obtido.Quantity.Equal(d("1.5")) {
		t.Fatalf("Buscar() = %+v, %v", obtido, err)
	}

	if err := repo.Remover(ctx, btc.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Buscar(ctx, btc.ID); !errors.Is(err, repository.ErrNaoEncontrado) {
		t.Fatalf("esperava ErrNaoEncontrado, obtido %v", err)
	}
	if err := repo.Remover(ctx, btc.ID); !errors.Is(err, repository.ErrNaoEncontrado) {
		t.Fatalf("esperava ErrNaoEncontrado ao remover de novo, obtido %v", err)
	}
	if err := repo.Atualizar(ctx, models.CryptoAsset{ID: 9999, Symbol: "X", Quantity: d("1"), Price: d("1")}); !errors.Is(err, repository.ErrNaoEncontrado) {
		t.Fatalf("esperava ErrNaoEncontrado ao atualizar inexistente, obtido %v", err)
	}
}
