# User Domain Documentation

The User Domain manages all aspects related to user accounts, including registration, authentication, profile management, and deletion.

## Key Components

*   **Entities (`domain/entities`, `internal/user/domain`)**:
    *   `User`: Represents a user account with attributes like ID, email, password hash, roles, etc.
    *   `RefreshToken`: Stores refresh tokens for maintaining user sessions.
    *   `UserDeletion`: Tracks user deletion requests and their status.
    *   `DeletionCapacity`: Manages the capacity for user deletions (e.g., rate limits).

*   **Repositories (`domain/repositories`, `internal/user/infrastructure`)**:
    *   Interfaces for persisting and retrieving user-related entities (e.g., `UserRepository`, `RefreshTokenRepository`).
    *   Implementations typically use PostgreSQL (`postgres_repository.go`) and Redis (`redis_blacklist_repository.go`, `redis_cache.go`).

*   **Application Services & Use Cases (`application/usecases`, `internal/user/application`)**:
    *   `Register`: Handles new user sign-ups, including password hashing and initial data storage.
    *   `Login`: Authenticates users, generates access and refresh tokens.
    *   `Logout`: Invalidates refresh tokens and blacklists access tokens.
    *   `PasswordRecovery`: Manages password reset requests and token verification.
    *   `DeleteUser`: Initiates and manages the user account deletion process.
    *   `VerifyToken`: Validates authentication tokens.

*   **Domain Services (`internal/user/domain`)**:
    *   `TokenService`: Handles the creation, signing, and verification of JWTs.
    *   `EventBus`: Publishes domain events (e.g., `UserRegistered`, `UserDeleted`) to the messaging system (NATS).

*   **Infrastructure (`internal/user/infrastructure`)**:
    *   `jwt_service.go`: Concrete implementation of `TokenService` using `golang-jwt`.
    *   `nats_event_bus.go`: Concrete implementation of `EventBus` using NATS.
    *   PostgreSQL repositories: Implementations for `UserRepository`, `RefreshTokenRepository`, etc.
    *   Redis components: `redis_blacklist_repository.go` for JWT blacklisting, `redis_cache.go` for general caching.

*   **Presentation (`internal/user/presentation`)**:
    *   `gin_handlers.go`: HTTP API endpoints for user-related operations using the Gin framework.
    *   `graphql_resolvers.go`: GraphQL resolvers for user-related queries and mutations.

## User Workflow Summary

```mermaid
graph TD
    A[User Request] --> B{Presentation Layer}
    B --> C[Application Layer (Use Cases)]
    C --> D{Domain Layer (Entities, Services)}
    D --> E[Infrastructure Layer (Repositories, External Services)]
    E --> F[Database/Redis/NATS]

    subgraph User Registration
        C1[Register Use Case] --> D1[User Entity Creation]
        D1 --> E1[UserRepository.Save]
        E1 --> F1[PostgreSQL]
        E1 --> F2[EventBus.Publish(UserRegistered)]
        F2 --> F3[NATS]
    end

    subgraph User Login
        C2[Login Use Case] --> D2[User Authentication]
        D2 --> D3[TokenService.GenerateTokens]
        D3 --> E2[RefreshTokenRepository.Save]
        E2 --> F1
    end

    subgraph User Deletion
        C3[DeleteUser Use Case] --> D4[UserDeletion Entity Creation]
        D4 --> E3[UserDeletionRepository.Save]
        E3 --> F1
        E3 --> F4[EventBus.Publish(UserDeletionRequested)]
        F4 --> F3
    end

    F3 --> G[Worker (e.g., User Delete Worker)]
```

**Bullet Summary of Core User Workflow:**

*   **Registration**: User provides credentials -> `Register` use case creates `User` entity, hashes password, saves to PostgreSQL, and publishes `UserRegistered` event to NATS.
*   **Login**: User provides credentials -> `Login` use case authenticates user, generates `AccessToken` and `RefreshToken` via `TokenService`, saves `RefreshToken` to PostgreSQL, and returns tokens.
*   **Authentication**: Incoming requests with `AccessToken` are validated by `TokenService`. If valid, the user's identity is established.
*   **Logout**: User requests logout -> `Logout` use case invalidates the `RefreshToken` and blacklists the `AccessToken` in Redis.
*   **Password Recovery**: User requests password reset -> `PasswordRecovery` use case generates a unique token, sends it via email, and allows password update upon token verification.
*   **User Deletion**: User requests account deletion -> `DeleteUser` use case records a `UserDeletion` entry, potentially checks `DeletionCapacity`, and publishes a `UserDeletionRequested` event to NATS for asynchronous processing by a worker.