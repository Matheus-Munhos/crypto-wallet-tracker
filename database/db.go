package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Conectar abre o pool de conexões, verifica a conexão e aplica as migrações.
func Conectar(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("nao foi possivel conectar ao banco: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("erro ao pingar o banco de dados: %w", err)
	}

	if err := Migrar(ctx, pool); err != nil {
		pool.Close()
		return nil, fmt.Errorf("erro ao criar tabela: %w", err)
	}
	return pool, nil
}

// Migrar garante que a tabela de ativos exista.
func Migrar(ctx context.Context, pool *pgxpool.Pool) error {
	const query = `
	CREATE TABLE IF NOT EXISTS crypto_assets (
		id SERIAL PRIMARY KEY,
		symbol VARCHAR(10) NOT NULL,
		quantity NUMERIC NOT NULL,
		price NUMERIC NOT NULL
	);`

	_, err := pool.Exec(ctx, query)
	return err
}
