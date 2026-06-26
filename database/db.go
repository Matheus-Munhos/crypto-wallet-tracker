package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConectarBanco() {

	connString := "postgres://admin:secret123@localhost:5432/cryptowallet"

	var err error

	DB, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Não foi possível conectar ao banco: %v\n", err)
		os.Exit(1)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		fmt.Printf("Erro ao pingar o banco de dados: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Conexão com o PostgreSQL do Docker feita com sucesso!")

	CriarTabela()
}

func CriarTabela() {
	query := `
	CREATE TABLE IF NOT EXISTS crypto_assets (
		id SERIAL PRIMARY KEY,
		symbol VARCHAR(10) NOT NULL,
		quantity NUMERIC NOT NULL,
		price NUMERIC NOT NULL
	);`

	_, err := DB.Exec(context.Background(), query)
	if err != nil {
		fmt.Printf("Erro ao criar tabela: %v\n", err)
		return
	}
	fmt.Println("Tabela 'crypto_assets' verificada/criada com sucesso!")
}
