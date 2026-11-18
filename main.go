package main

import (
	"context"
	"crypto-tracker/config"
	"crypto-tracker/controllers"
	"crypto-tracker/database"
	"crypto-tracker/repository"
	"crypto-tracker/routes"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Carregar()
	if err != nil {
		log.Fatal(err)
	}

	// ctx é cancelado quando o processo recebe Ctrl+C ou SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.Conectar(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	log.Println("Conexao com o PostgreSQL feita com sucesso!")

	assetHandler := controllers.NovoAssetHandler(repository.NovoPostgresRepository(pool))

	servidor := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           routes.ConfigurarRotas(assetHandler),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	erros := make(chan error, 1)
	go func() {
		log.Printf("Servidor da Carteira Crypto rodando na porta :%s...", cfg.Port)
		if err := servidor.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			erros <- err
		}
	}()

	select {
	case err := <-erros:
		log.Printf("erro ao iniciar o servidor: %v", err)
		return
	case <-ctx.Done():
	}

	log.Println("Encerrando o servidor...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := servidor.Shutdown(shutdownCtx); err != nil {
		log.Printf("erro ao encerrar o servidor: %v", err)
	}
}
