FROM golang:1.25-alpine
WORKDIR /app

# Pre-caché de módulos
COPY go.mod go.sum ./
RUN go mod download

# Para ver logs con timestamps legibles
ENV GIN_MODE=debug

# El código se monta vía volumen en docker-compose.yml
CMD ["go", "run", "cmd/main.go"]
