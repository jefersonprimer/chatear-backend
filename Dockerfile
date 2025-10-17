# --- Stage 1: Build ---
FROM golang:1.21-alpine AS builder

# Instala dependências necessárias
RUN apk add --no-cache git ca-certificates build-base

# Define diretório de trabalho
WORKDIR /app

# Copia módulos e baixa dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código
COPY . .

# Builda o binário
RUN go build -o chatear-backend ./cmd

# --- Stage 2: Production ---
FROM alpine:3.18

# Instala certificado de SSL (necessário para conexões externas)
RUN apk add --no-cache ca-certificates

# Cria diretório do app
WORKDIR /app

# Copia binário do stage anterior
COPY --from=builder /app/chatear-backend .

# Expõe a porta que seu docker-compose mapeia
EXPOSE 8080

# Comando para rodar seu app
CMD ["./chatear-backend"]

