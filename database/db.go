package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// ConectarBanco abre o pool de conexões, verifica a conexão e garante que a tabela exista.
func ConectarBanco(ctx context.Context, connString string) error {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return fmt.Errorf("nao foi possivel conectar ao banco: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return fmt.Errorf("erro ao pingar o banco de dados: %w", err)
	}

	DB = pool
	log.Println("Conexao com o PostgreSQL feita com sucesso!")

	if err := CriarTabela(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("erro ao criar tabela: %w", err)
	}
	return nil
}

func CriarTabela(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS crypto_assets (
		id SERIAL PRIMARY KEY,
		symbol VARCHAR(10) NOT NULL,
		quantity NUMERIC NOT NULL,
		price NUMERIC NOT NULL
	);`

	if _, err := DB.Exec(ctx, query); err != nil {
		return err
	}
	log.Println("Tabela 'crypto_assets' verificada/criada com sucesso!")
	return nil
}
