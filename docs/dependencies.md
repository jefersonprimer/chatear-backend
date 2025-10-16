# Project Dependencies

This document outlines the third-party libraries used in the project and their primary purposes.

## Direct Dependencies

*   **`github.com/99designs/gqlgen`** (GraphQL Server Library):
    *   **Purpose**: Used for building GraphQL servers in Go. It generates a GraphQL schema from Go types and provides a runtime for executing GraphQL queries.

*   **`github.com/gin-gonic/gin`** (Web Framework):
    *   **Purpose**: A high-performance HTTP web framework for Go. It's used for handling HTTP requests, routing, and middleware in the API server.

*   **`github.com/go-redis/redis/v8`** and **`github.com/redis/go-redis/v9`** (Redis Clients):
    *   **Purpose**: Clients for interacting with Redis, an in-memory data structure store. Used for caching, session management, and potentially rate limiting.

*   **`github.com/golang-jwt/jwt/v4`** and **`github.com/golang-jwt/jwt/v5`** (JWT Implementation):
    *   **Purpose**: Libraries for working with JSON Web Tokens (JWTs), used for authentication and authorization within the application.

*   **`github.com/google/uuid`** (UUID Generation):
    *   **Purpose**: Provides functionality to generate and parse Universally Unique Identifiers (UUIDs).

*   **`github.com/jackc/pgx/v5`** (PostgreSQL Driver):
    *   **Purpose**: A high-performance PostgreSQL driver and toolkit for Go, used for interacting with the PostgreSQL database.

*   **`github.com/joho/godotenv`** (Environment Variable Loader):
    *   **Purpose**: Loads environment variables from `.env` files into `os.Getenv`, simplifying local development configuration.

*   **`github.com/nats-io/nats.go`** (NATS Client):
    *   **Purpose**: The official Go client for NATS, a high-performance messaging system. Used for inter-service communication and event streaming (e.g., for workers).

*   **`github.com/stretchr/testify`** (Testing Toolkit):
    *   **Purpose**: Provides a set of useful tools for writing tests in Go, including assertion functions and mocking capabilities.

*   **`github.com/vektah/gqlparser/v2`** (GraphQL Parser):
    *   **Purpose**: A GraphQL parser and validator, often used internally by `gqlgen` to process GraphQL schema definitions and queries.

*   **`golang.org/x/crypto`** (Cryptographic Functions):
    *   **Purpose**: Provides various cryptographic functionalities, such as password hashing (e.g., bcrypt) and other security-related operations.

## Indirect Dependencies

Numerous other indirect dependencies are pulled in by the direct dependencies. These are managed by Go Modules and are typically not directly interacted with by the application code but are essential for the direct dependencies to function correctly. They are listed in `go.mod` under the `require` block with `// indirect` comments.
