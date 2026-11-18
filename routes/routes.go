package routes

import (
	"crypto-tracker/controllers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ConfigurarRotas(assets *controllers.AssetHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	r.Route("/wallet", func(router chi.Router) {
		router.Post("/", assets.CadastrarAtivo)
		router.Get("/", assets.ListarAtivos)
		router.Get("/total", assets.CalcularTotal)
		router.Get("/{id}", assets.BuscarAtivo)
		router.Put("/{id}", assets.AtualizarAtivo)
		router.Delete("/{id}", assets.RemoverAtivo)
	})

	return r
}
