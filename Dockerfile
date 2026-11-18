# --- build ---
FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/crypto-tracker .

# --- runtime ---
# distroless/static já inclui os certificados CA (necessários para chamar a CoinGecko via HTTPS)
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/crypto-tracker /crypto-tracker

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/crypto-tracker"]
