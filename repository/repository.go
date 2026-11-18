package repository

import (
	"context"
	"crypto-tracker/models"
	"errors"
)

// ErrNaoEncontrado é retornado quando o ativo pedido não existe.
var ErrNaoEncontrado = errors.New("ativo nao encontrado")

// AssetRepository define as operações de persistência dos ativos.
type AssetRepository interface {
	Criar(ctx context.Context, asset *models.CryptoAsset) error
	Listar(ctx context.Context) ([]models.CryptoAsset, error)
	Buscar(ctx context.Context, id int) (models.CryptoAsset, error)
	Atualizar(ctx context.Context, asset models.CryptoAsset) error
	Remover(ctx context.Context, id int) error
	TotalInvestido(ctx context.Context) (float64, error)
}
