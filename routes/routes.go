package routes

import (
	controller "crypto-tracker/controllers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ConfigurarRotas() http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write(([]byte("pong")))
	})

	r.Route("/wallet", func(router chi.Router) {
		router.Post("/", controller.CadastrarAtivo)
		router.Get("/", controller.ListarArquivos)
		router.Get("/total", controller.CalcularTotal)
	})

	return r

}
