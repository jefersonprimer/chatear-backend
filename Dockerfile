# ===============================
# Stage 1: Build
# ===============================
FROM golang:1.25.3-alpine AS builder

# Instala dependências necessárias
RUN apk add --no-cache git ca-certificates build-base

# Define diretório de trabalho
WORKDIR /app

# Copia módulos e baixa dependências
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copia todo o código
COPY . .

# Builda os binários
RUN go build -o chatear-api ./cmd/api
RUN go build -o chatear-worker ./cmd/worker

# ===============================
# Stage 2: Production
# ===============================
FROM alpine:3.18

# Instala certificado de SSL
RUN apk add --no-cache ca-certificates

# Cria diretório do app
WORKDIR /app

# Copia binários do stage de build
COPY --from=builder /app/chatear-api .
COPY --from=builder /app/chatear-worker .

# Expõe porta da API
EXPOSE 8080

# Define variável de ambiente para escolha do binário (API ou Worker)
ENV APP_BIN=chatear-api

# Comando padrão (API). Para rodar o worker, sobrescreve no docker-compose ou CLI
CMD ["sh", "-c", "./$APP_BIN"]

