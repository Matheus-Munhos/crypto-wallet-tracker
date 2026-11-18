package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

// tamanhoMaximoCorpo limita o corpo das requisições a 1 MB.
const tamanhoMaximoCorpo = 1 << 20

func lerJSON(w http.ResponseWriter, r *http.Request, destino any) error {
	r.Body = http.MaxBytesReader(w, r.Body, tamanhoMaximoCorpo)
	return json.NewDecoder(r.Body).Decode(destino)
}

func escreverJSON(w http.ResponseWriter, status int, dados any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(dados); err != nil {
		log.Printf("erro ao escrever resposta JSON: %v", err)
	}
}

func escreverErro(w http.ResponseWriter, status int, mensagem string) {
	escreverJSON(w, status, map[string]string{"error": mensagem})
}
