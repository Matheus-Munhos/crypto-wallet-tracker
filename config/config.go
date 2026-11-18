package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config reúne as configurações da aplicação lidas do ambiente.
type Config struct {
	DatabaseURL string
	Port        string
}

// Carregar lê as variáveis de ambiente (e o arquivo .env, se existir).
func Carregar() (Config, error) {
	// O .env é opcional: em produção as variáveis vêm do próprio ambiente.
	_ = godotenv.Load()

	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        valorOuPadrao("PORT", "8080"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("variavel de ambiente DATABASE_URL nao definida (veja o .env.example)")
	}

	return cfg, nil
}

func valorOuPadrao(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}
