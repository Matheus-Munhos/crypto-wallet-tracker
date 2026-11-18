.PHONY: run build test test-integration vet fmt db-up up down logs

## run: roda a API localmente (lê o .env)
run:
	go run .

## build: gera o binário em bin/
build:
	go build -o bin/crypto-tracker .

## test: roda os testes unitários
test:
	go test ./...

## test-integration: roda também os testes com PostgreSQL (requer TEST_DATABASE_URL)
test-integration:
	@test -n "$(TEST_DATABASE_URL)" || (echo "defina TEST_DATABASE_URL" && exit 1)
	go test -count=1 ./...

## vet: análise estática
vet:
	go vet ./...

## fmt: formata o código
fmt:
	gofmt -w .

## db-up: sobe apenas o PostgreSQL
db-up:
	docker compose up -d postgres

## up: sobe PostgreSQL + API com Docker
up:
	docker compose up -d --build

## down: derruba os containers
down:
	docker compose down

## logs: acompanha os logs da API
logs:
	docker compose logs -f api
