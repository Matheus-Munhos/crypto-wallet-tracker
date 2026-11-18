package main

import (
	"crypto-tracker/database"
	"crypto-tracker/routes"
	"log"
	"net/http"
)

func main() {

	database.ConectarBanco()
	roteador := routes.ConfigurarRotas()

	log.Println("Servidor da Carteira Crypto rodando na porta :8080...")

	if err := http.ListenAndServe(":8080", roteador); err != nil {
		log.Fatalf("erro ao iniciar o servidor: %v", err)
	}
}
