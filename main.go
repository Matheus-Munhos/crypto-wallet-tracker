package main

import (
	"crypto-tracker/database"
	"crypto-tracker/routes"
	"fmt"
	"net/http"
)

func main() {

	database.ConectarBanco()
	roteador := routes.ConfigurarRotas()

	fmt.Println("Servidor da Carteira Crypto rodando na porta :8080...")

	http.ListenAndServe(":8080", roteador)

}
