# Guia de Início Rápido para o Projeto Chatear-Backend

Este guia detalha como configurar, executar e testar o projeto Chatear-Backend.

## 1. Pré-requisitos

Certifique-se de ter as seguintes ferramentas instaladas em sua máquina:

*   **Go:** Versão 1.21 ou superior.
*   **Docker:** Para executar os serviços de banco de dados, cache e mensageria.
*   **Docker Compose:** Para orquestrar os contêineres Docker.
*   **Insomnia (ou Postman):** Para testar as APIs HTTP e GraphQL.

## 2. Configuração do Projeto

1.  **Clonar o Repositório:**
    ```bash
    git clone https://github.com/jefersonprimer/chatear-backend.git
    cd chatear-backend
    ```

2.  **Instalar Dependências Go:**
    ```bash
    go mod download
    ```

3.  **Configurar Variáveis de Ambiente:**
    Crie um arquivo `.env` na raiz do projeto, copiando o conteúdo de `env.example` e preenchendo as variáveis conforme necessário.

    ```bash
    cp env.example .env
    ```

    Exemplo de `.env` (preencha com seus valores):
    ```
    # Application
    APP_ENV=development
    APP_URL=http://localhost:8080
    JWT_SECRET=your_jwt_secret_key
    REFRESH_TOKEN_SECRET=your_refresh_token_secret_key

    # Database (Postgres)
    DATABASE_URL=postgres://user:password@localhost:5432/chatear_db?sslmode=disable

    # Redis
    REDIS_URL=redis://localhost:6379/0

    # NATS
    NATS_URL=nats://localhost:4222

    # Supabase (optional, if used)
    SUPABASE_URL=
    SUPABASE_ANON_KEY=

    # SMTP for Email Sending
    SMTP_HOST=smtp.mailtrap.io
    SMTP_PORT=2525
    SMTP_USER=your_mailtrap_username
    SMTP_PASS=your_mailtrap_password
    SMTP_FROM=no-reply@chatear.com
    ```

## 3. Executando os Serviços Docker

Os serviços de banco de dados (Postgres), cache (Redis) e mensageria (NATS) são essenciais para o funcionamento do backend.

1.  **Iniciar os Contêineres:**
    ```bash
    docker-compose -f docker-compose.events.yml up -d
    ```
    Isso iniciará o Postgres, Redis e NATS em segundo plano.

2.  **Verificar o Status:**
    ```bash
    docker-compose -f docker-compose.events.yml ps
    ```

## 4. Executando a Aplicação Go

O backend é composto por um servidor API principal e workers para processamento assíncrono.

1.  **Executar o Servidor API:**
    Abra um novo terminal e execute:
    ```bash
    go run cmd/api/main.go
    ```
    O servidor API estará disponível em `http://localhost:8080` (ou a porta configurada).

2.  **Executar o Worker de Notificação (Email):**
    Abra outro terminal e execute:
    ```bash
    go run cmd/worker/notification_worker.go
    ```
    Este worker processará os eventos de envio de e-mail publicados via NATS.

3.  **Executar o Worker de Exclusão de Usuário:**
    Abra outro terminal e execute:
    ```bash
    go run cmd/worker/user_delete_worker.go
    ```
    Este worker processará as solicitações de exclusão de usuário.

## 5. Testando com Insomnia

Você pode testar as APIs HTTP e GraphQL usando o Insomnia.

### 5.1. Importar Coleção (se disponível)

Se houver um arquivo de coleção do Insomnia (`.json` ou `.yaml`) no diretório `docs/`, importe-o diretamente. Caso contrário, siga os passos abaixo para criar as requisições manualmente.

### 5.2. Requisições HTTP (Exemplos)

**Base URL:** `http://localhost:8080`

1.  **Registro de Usuário:**
    *   **Método:** `POST`
    *   **URL:** `/register`
    *   **Body (JSON):**
        ```json
        {
            "name": "Test User",
            "email": "test@example.com",
            "password": "password123"
        }
        ```

2.  **Login de Usuário:**
    *   **Método:** `POST`
    *   **URL:** `/login`
    *   **Body (JSON):**
        ```json
        {
            "email": "test@example.com",
            "password": "password123"
        }
        ```
    *   **Observação:** Guarde o `access_token` e `refresh_token` da resposta.

3.  **Verificar E-mail:**
    *   **Método:** `GET`
    *   **URL:** `/verify-email?token=<TOKEN_DE_VERIFICACAO>`
    *   **Observação:** O token de verificação será enviado para o e-mail registrado (verifique os logs do worker de notificação ou o serviço SMTP configurado, como Mailtrap).

4.  **Obter Perfil do Usuário (`/me`):**
    *   **Método:** `GET`
    *   **URL:** `/me`
    *   **Headers:**
        *   `Authorization`: `Bearer <SEU_ACCESS_TOKEN>`

5.  **Logout de Usuário:**
    *   **Método:** `POST`
    *   **URL:** `/logout`
    *   **Headers:**
        *   `Authorization`: `Bearer <SEU_ACCESS_TOKEN>`
    *   **Body (JSON):**
        ```json
        {
            "refresh_token": "<SEU_REFRESH_TOKEN>"
        }
        ```

### 5.3. Requisições GraphQL (Exemplos)

**Endpoint:** `http://localhost:8080/query`

1.  **Registro de Usuário (Mutation `registerUser`):**
    *   **Método:** `POST`
    *   **URL:** `/query`
    *   **Body (GraphQL Query):**
        ```graphql
        mutation RegisterUser {
          registerUser(input: {name: "GraphQL User", email: "graphql@example.com", password: "password123"}) {
            user {
              id
              name
              email
              isEmailVerified
            }
            accessToken
            refreshToken
          }
        }
        ```

2.  **Login de Usuário (Mutation `login`):**
    *   **Método:** `POST`
    *   **URL:** `/query`
    *   **Body (GraphQL Query):**
        ```graphql
        mutation Login {
          login(input: {email: "graphql@example.com", password: "password123"}) {
            user {
              id
              email
            }
            accessToken
            refreshToken
          }
        }
        ```
    *   **Observação:** Guarde o `accessToken` e `refreshToken` da resposta.

3.  **Obter Perfil do Usuário (Query `me`):**
    *   **Método:** `POST`
    *   **URL:** `/query`
    *   **Headers:**
        *   `Authorization`: `Bearer <SEU_ACCESS_TOKEN>`
    *   **Body (GraphQL Query):**
        ```graphql
        query Me {
          me {
            id
            name
            email
            isEmailVerified
          }
        }
        ```

4.  **Excluir Conta (Mutation `deleteAccount`):**
    *   **Método:** `POST`
    *   **URL:** `/query`
    *   **Headers:**
        *   `Authorization`: `Bearer <SEU_ACCESS_TOKEN>`
    *   **Body (GraphQL Query):**
        ```graphql
        mutation DeleteAccount {
          deleteAccount(input: {userID: "<ID_DO_USUARIO_A_EXCLUIR>"})
        }
        ```
    *   **Observação:** O `userID` pode ser obtido após o registro ou login.

## 6. Verificação

*   **Logs:** Monitore os terminais onde o servidor API e os workers estão sendo executados para verificar mensagens de log e erros.
*   **Banco de Dados:** Conecte-se ao banco de dados Postgres (por exemplo, usando `psql` ou uma ferramenta GUI) para verificar se os usuários estão sendo criados e atualizados corretamente.
*   **Redis:** Use um cliente Redis para verificar se os tokens estão sendo armazenados.
*   **Mailtrap (ou serviço SMTP configurado):** Verifique a caixa de entrada para e-mails de verificação ou outros e-mails transacionais.

Este guia deve fornecer um ponto de partida sólido para explorar e testar o Chatear-Backend.
