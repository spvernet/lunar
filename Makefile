APP_NAME=lunar-rockets
DC=docker-compose

.PHONY: tidy fmt build run dev test stop


install:
	@$(DC) run --rm app go mod download #@go mod download

tidy:
	@go mod tidy

fmt:
	@go fmt ./...

# Arranque dev con volumen + go run
run:
	@$(DC) up --build


test:
	@$(DC) run --rm app go test ./...

stop:
	@$(DC) down --remove-orphans
