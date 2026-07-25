.PHONY: help build run test test-race test-cover lint vet fmt vuln migrate-up migrate-down mocks clean

APP_NAME := api
BIN_DIR  := bin
MAIN_PATH := ./cmd/api

help: ## Mostra os comandos disponíveis
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	
build: ## Compila o binário
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PATH)
	
run: ## Roda a aplicação localmente
	go run $(MAIN_PATH)

test: ## Roda os testes unitários
	go test ./... -v

test-race: ## Roda testes com detector de data race
	go test -race -count=1 ./...

test-cover: ## Roda testes com cobertura
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	
lint: ## Roda o linter
	golangci-lint run ./...

vet: ## Roda go vet
	go vet ./...

fmt: ## Formata o código (gofumpt)
	gofumpt -l -w .

vuln: ## Verifica vulnerabilidades em dependências
	govulncheck ./...

secrets: ## Escaneia segredos vazados (gitleaks)
	gitleaks detect --source . -v

mocks: ## Gera mocks (mockery)
	mockery --all --output ./internal/mocks
	
check: fmt vet lint test-race vuln ## Roda tudo que o CI roda, localmente
	
clean: ## Remove binários e artefatos
	rm -rf $(BIN_DIR) coverage.out