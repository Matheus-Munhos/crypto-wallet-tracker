package repository

import (
	"context"
	"crypto-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository implementa AssetRepository usando PostgreSQL.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NovoPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Criar(ctx context.Context, asset *models.CryptoAsset) error {
	const query = `
		INSERT INTO crypto_assets (symbol, quantity, price)
		VALUES ($1, $2, $3)
		RETURNING id`
	return r.pool.QueryRow(ctx, query, asset.Symbol, asset.Quantity, asset.Price).Scan(&asset.ID)
}

func (r *PostgresRepository) Listar(ctx context.Context) ([]models.CryptoAsset, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, symbol, quantity, price FROM crypto_assets ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assets := []models.CryptoAsset{}
	for rows.Next() {
		var a models.CryptoAsset
		if err := rows.Scan(&a.ID, &a.Symbol, &a.Quantity, &a.Price); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

func (r *PostgresRepository) TotalInvestido(ctx context.Context) (float64, error) {
	var total float64
	err := r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(quantity * price), 0) FROM crypto_assets`).Scan(&total)
	return total, err
}
