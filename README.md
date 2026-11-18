<a id="english"></a>

# Crypto Wallet Tracker

[![CI](https://github.com/Matheus-Munhos/crypto-wallet-tracker/actions/workflows/ci.yml/badge.svg)](https://github.com/Matheus-Munhos/crypto-wallet-tracker/actions/workflows/ci.yml)

> 🇧🇷 **Há uma versão deste README em português mais abaixo. [Clique aqui para ler em português](#portugues).**

A REST API written in Go to track a crypto wallet: register the assets you bought, see how much you invested and compare it with the current market value (prices from [CoinGecko](https://www.coingecko.com/)).

## Features

- Full CRUD for wallet assets (symbol, quantity and purchase price)
- Total invested in the wallet
- Wallet valuation at current prices: current value and profit/loss (absolute and %), in total and per asset
- Exact decimal arithmetic (no floating-point rounding errors)
- Input validation with clear JSON error messages
- Configuration through environment variables / `.env`
- Graceful shutdown, server timeouts and request logging
- Docker image (distroless, ~22 MB) and `docker compose` setup
- Unit and integration tests, plus CI with GitHub Actions

## Tech stack

- [Go](https://go.dev/) 1.26+
- [chi](https://github.com/go-chi/chi) (router)
- [pgx](https://github.com/jackc/pgx) (PostgreSQL driver)
- [shopspring/decimal](https://github.com/shopspring/decimal) (decimal numbers)
- PostgreSQL 15
- Docker / Docker Compose

## Getting started

### 1. Configure the environment

```bash
cp .env.example .env        # PowerShell: Copy-Item .env.example .env
```

Edit `.env` and set at least `POSTGRES_PASSWORD` (and the same password inside `DATABASE_URL`).

### 2a. Run everything with Docker

```bash
docker compose up -d --build
```

The API will be available at `http://localhost:8080` (change it with `API_PORT`).

### 2b. Or run the API locally

```bash
docker compose up -d postgres   # only the database
go run .
```

> If port `5432` is already in use on your machine, set `POSTGRES_PORT` in `.env` (e.g. `5433`) and update the port in `DATABASE_URL`.

A `Makefile` with shortcuts is also available (`make run`, `make test`, `make up`, `make down`...).

## Configuration

| Variable            | Default                            | Description                                                                  |
| ------------------- | ---------------------------------- | ---------------------------------------------------------------------------- |
| `DATABASE_URL`      | — (required)                       | PostgreSQL connection string used by the API                                 |
| `PORT`              | `8080`                             | Port the API listens on                                                      |
| `PRICE_CURRENCY`    | `usd`                              | Currency for current prices (`usd`, `brl`, `eur`...). **Must match the currency of the prices you register** |
| `COINGECKO_API_URL` | `https://api.coingecko.com/api/v3` | CoinGecko API base URL                                                       |
| `COINGECKO_API_KEY` | —                                  | Optional CoinGecko demo API key (higher rate limits)                         |
| `POSTGRES_USER`     | `admin`                            | Database user (docker compose)                                               |
| `POSTGRES_PASSWORD` | — (required)                       | Database password (docker compose)                                           |
| `POSTGRES_DB`       | `cryptowallet`                     | Database name (docker compose)                                               |
| `POSTGRES_PORT`     | `5432`                             | Database port on the host (docker compose)                                   |
| `API_PORT`          | `8080`                             | API port on the host (docker compose)                                        |

## API

| Method   | Endpoint            | Description                              |
| -------- | ------------------- | ---------------------------------------- |
| `GET`    | `/ping`             | Health check (returns `pong`)            |
| `POST`   | `/wallet`           | Register an asset                        |
| `GET`    | `/wallet`           | List all assets                          |
| `GET`    | `/wallet/{id}`      | Get an asset by id                       |
| `PUT`    | `/wallet/{id}`      | Update an asset                          |
| `DELETE` | `/wallet/{id}`      | Delete an asset                          |
| `GET`    | `/wallet/total`     | Total invested                           |
| `GET`    | `/wallet/valuation` | Valuation at current prices + profit/loss |

> Monetary values (`quantity`, `price`, totals) are returned as **strings** to preserve precision (e.g. `"0.1"`). As input, both numbers and strings are accepted.

### Register an asset

```bash
curl -X POST http://localhost:8080/wallet \
  -H "Content-Type: application/json" \
  -d '{"symbol": "btc", "quantity": "0.5", "price": "60000"}'
```

```json
{ "id": 1, "symbol": "BTC", "quantity": "0.5", "price": "60000" }
```

Validation rules: `symbol` must have 1–10 alphanumeric characters (it is converted to uppercase), `quantity` must be greater than zero and `price` cannot be negative.

### Total invested

```bash
curl http://localhost:8080/wallet/total
```

```json
{ "total_invested": "30000" }
```

### Valuation at current prices

```bash
curl http://localhost:8080/wallet/valuation
```

```json
{
  "currency": "usd",
  "total_invested": "30000",
  "current_value": "31393.5",
  "profit_loss": "1393.5",
  "profit_loss_percent": "4.65",
  "assets": [
    {
      "id": 1,
      "symbol": "BTC",
      "quantity": "0.5",
      "price": "60000",
      "current_price": "62787",
      "invested": "30000",
      "current_value": "31393.5",
      "profit_loss": "1393.5",
      "profit_loss_percent": "4.65"
    }
  ],
  "missing_prices": []
}
```

Symbols without a price on CoinGecko are listed in `missing_prices` and are left out of the totals. If CoinGecko is unavailable, the endpoint returns `502 Bad Gateway`.

### Errors

Errors are always JSON:

```json
{ "error": "quantity deve ser maior que zero" }
```

| Status | When                                     |
| ------ | ---------------------------------------- |
| `400`  | Invalid JSON, invalid fields or invalid id |
| `404`  | Asset not found                          |
| `500`  | Database error                           |
| `502`  | Price provider (CoinGecko) unavailable   |

## Tests

```bash
go test ./...
```

The PostgreSQL integration test only runs when `TEST_DATABASE_URL` is set. **Warning:** it truncates the `crypto_assets` table, so use a dedicated database:

```bash
TEST_DATABASE_URL="postgres://admin:change-me@localhost:5432/cryptowallet_test?sslmode=disable" go test -count=1 ./...
```

## Project structure

```
.
├── config/        # Environment configuration
├── controllers/   # HTTP handlers
├── database/      # Connection and migration
├── models/        # Entities, validation and valuation logic
├── pricing/       # Current price provider (CoinGecko)
├── repository/    # Data access (PostgreSQL)
├── routes/        # Routes and middlewares
└── main.go        # Entry point
```

---

<a id="portugues"></a>

# Crypto Wallet Tracker (Português)

> 🇺🇸 [Back to the English version](#english)

API REST em Go para acompanhar uma carteira de criptomoedas: cadastre os ativos que você comprou, veja quanto investiu e compare com o valor atual de mercado (cotações da [CoinGecko](https://www.coingecko.com/)).

## Funcionalidades

- CRUD completo dos ativos da carteira (símbolo, quantidade e preço de compra)
- Total investido na carteira
- Avaliação da carteira a preços atuais: valor atual e lucro/prejuízo (absoluto e %), no total e por ativo
- Aritmética decimal exata (sem erros de arredondamento de ponto flutuante)
- Validação de entrada com mensagens de erro claras em JSON
- Configuração por variáveis de ambiente / `.env`
- Encerramento limpo (graceful shutdown), timeouts no servidor e log das requisições
- Imagem Docker (distroless, ~22 MB) e configuração com `docker compose`
- Testes unitários e de integração, além de CI com GitHub Actions

## Tecnologias

- [Go](https://go.dev/) 1.26+
- [chi](https://github.com/go-chi/chi) (roteador)
- [pgx](https://github.com/jackc/pgx) (driver do PostgreSQL)
- [shopspring/decimal](https://github.com/shopspring/decimal) (números decimais)
- PostgreSQL 15
- Docker / Docker Compose

## Como rodar

### 1. Configure o ambiente

```bash
cp .env.example .env        # PowerShell: Copy-Item .env.example .env
```

Edite o `.env` e defina pelo menos `POSTGRES_PASSWORD` (e a mesma senha dentro de `DATABASE_URL`).

### 2a. Rodar tudo com Docker

```bash
docker compose up -d --build
```

A API fica disponível em `http://localhost:8080` (altere com `API_PORT`).

### 2b. Ou rodar a API localmente

```bash
docker compose up -d postgres   # apenas o banco
go run .
```

> Se a porta `5432` já estiver em uso na sua máquina, defina `POSTGRES_PORT` no `.env` (ex.: `5433`) e ajuste a porta na `DATABASE_URL`.

Também há um `Makefile` com atalhos (`make run`, `make test`, `make up`, `make down`...).

## Configuração

| Variável            | Padrão                             | Descrição                                                                    |
| ------------------- | ---------------------------------- | ---------------------------------------------------------------------------- |
| `DATABASE_URL`      | — (obrigatória)                    | String de conexão do PostgreSQL usada pela API                               |
| `PORT`              | `8080`                             | Porta em que a API escuta                                                    |
| `PRICE_CURRENCY`    | `usd`                              | Moeda das cotações atuais (`usd`, `brl`, `eur`...). **Deve ser a mesma moeda dos preços cadastrados** |
| `COINGECKO_API_URL` | `https://api.coingecko.com/api/v3` | URL base da API da CoinGecko                                                 |
| `COINGECKO_API_KEY` | —                                  | Chave demo opcional da CoinGecko (limites de requisição maiores)             |
| `POSTGRES_USER`     | `admin`                            | Usuário do banco (docker compose)                                            |
| `POSTGRES_PASSWORD` | — (obrigatória)                    | Senha do banco (docker compose)                                              |
| `POSTGRES_DB`       | `cryptowallet`                     | Nome do banco (docker compose)                                               |
| `POSTGRES_PORT`     | `5432`                             | Porta do banco no host (docker compose)                                      |
| `API_PORT`          | `8080`                             | Porta da API no host (docker compose)                                        |

## API

| Método   | Endpoint            | Descrição                                   |
| -------- | ------------------- | ------------------------------------------- |
| `GET`    | `/ping`             | Verificação de saúde (retorna `pong`)       |
| `POST`   | `/wallet`           | Cadastra um ativo                           |
| `GET`    | `/wallet`           | Lista todos os ativos                       |
| `GET`    | `/wallet/{id}`      | Busca um ativo pelo id                      |
| `PUT`    | `/wallet/{id}`      | Atualiza um ativo                           |
| `DELETE` | `/wallet/{id}`      | Remove um ativo                             |
| `GET`    | `/wallet/total`     | Total investido                             |
| `GET`    | `/wallet/valuation` | Avaliação a preços atuais + lucro/prejuízo  |

> Valores monetários (`quantity`, `price`, totais) são retornados como **string** para preservar a precisão (ex.: `"0.1"`). Na entrada, são aceitos números ou strings.

### Cadastrar um ativo

```bash
curl -X POST http://localhost:8080/wallet \
  -H "Content-Type: application/json" \
  -d '{"symbol": "btc", "quantity": "0.5", "price": "60000"}'
```

```json
{ "id": 1, "symbol": "BTC", "quantity": "0.5", "price": "60000" }
```

Regras de validação: `symbol` deve ter de 1 a 10 caracteres alfanuméricos (é convertido para maiúsculas), `quantity` deve ser maior que zero e `price` não pode ser negativo.

### Total investido

```bash
curl http://localhost:8080/wallet/total
```

```json
{ "total_invested": "30000" }
```

### Avaliação a preços atuais

```bash
curl http://localhost:8080/wallet/valuation
```

A resposta tem o mesmo formato do exemplo da [versão em inglês](#valuation-at-current-prices): totais (`total_invested`, `current_value`, `profit_loss`, `profit_loss_percent`), a lista `assets` com o resultado de cada ativo e `missing_prices`.

Símbolos sem cotação na CoinGecko aparecem em `missing_prices` e ficam fora dos totais. Se a CoinGecko estiver indisponível, o endpoint retorna `502 Bad Gateway`.

### Erros

Os erros são sempre em JSON:

```json
{ "error": "quantity deve ser maior que zero" }
```

| Status | Quando                                       |
| ------ | -------------------------------------------- |
| `400`  | JSON inválido, campos inválidos ou id inválido |
| `404`  | Ativo não encontrado                         |
| `500`  | Erro no banco de dados                       |
| `502`  | Provedor de cotações (CoinGecko) indisponível |

## Testes

```bash
go test ./...
```

O teste de integração com PostgreSQL só roda quando `TEST_DATABASE_URL` está definida. **Atenção:** ele apaga a tabela `crypto_assets`, então use um banco dedicado:

```bash
TEST_DATABASE_URL="postgres://admin:change-me@localhost:5432/cryptowallet_test?sslmode=disable" go test -count=1 ./...
```

## Estrutura do projeto

```
.
├── config/        # Configuração por variáveis de ambiente
├── controllers/   # Handlers HTTP
├── database/      # Conexão e migração
├── models/        # Entidades, validação e cálculo da avaliação
├── pricing/       # Provedor de cotações atuais (CoinGecko)
├── repository/    # Acesso a dados (PostgreSQL)
├── routes/        # Rotas e middlewares
└── main.go        # Ponto de entrada
```
