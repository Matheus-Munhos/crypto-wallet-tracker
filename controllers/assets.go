package controllers

import (
	"crypto-tracker/models"
	"crypto-tracker/repository"
	"log"
	"net/http"
)

// AssetHandler agrupa os handlers HTTP da carteira.
type AssetHandler struct {
	repo repository.AssetRepository
}

func NovoAssetHandler(repo repository.AssetRepository) *AssetHandler {
	return &AssetHandler{repo: repo}
}

// CadastrarAtivo trata POST /wallet.
func (h *AssetHandler) CadastrarAtivo(w http.ResponseWriter, r *http.Request) {
	var asset models.CryptoAsset
	if err := lerJSON(w, r, &asset); err != nil {
		escreverErro(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	asset.Normalizar()
	if err := asset.Validar(); err != nil {
		escreverErro(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Criar(r.Context(), &asset); err != nil {
		log.Printf("erro ao salvar ativo: %v", err)
		escreverErro(w, http.StatusInternalServerError, "Erro ao salvar no banco de dados")
		return
	}

	escreverJSON(w, http.StatusCreated, asset)
}

// ListarAtivos trata GET /wallet.
func (h *AssetHandler) ListarAtivos(w http.ResponseWriter, r *http.Request) {
	assets, err := h.repo.Listar(r.Context())
	if err != nil {
		log.Printf("erro ao listar ativos: %v", err)
		escreverErro(w, http.StatusInternalServerError, "Erro ao listar ativos")
		return
	}
	escreverJSON(w, http.StatusOK, assets)
}

// CalcularTotal trata GET /wallet/total.
func (h *AssetHandler) CalcularTotal(w http.ResponseWriter, r *http.Request) {
	total, err := h.repo.TotalInvestido(r.Context())
	if err != nil {
		log.Printf("erro ao calcular total: %v", err)
		escreverErro(w, http.StatusInternalServerError, "Erro ao calcular total")
		return
	}
	escreverJSON(w, http.StatusOK, map[string]float64{"total_invested": total})
}
