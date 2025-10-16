# Architecture Overview

This project follows a layered architectural pattern, aiming for separation of concerns, maintainability, and testability. The primary layers are:

1.  **Presentation Layer (`cmd/api`, `graph`, `presentation`)**:
    *   Handles incoming requests (HTTP, GraphQL).
    *   Translates requests into commands or queries for the Application Layer.
    *   Serializes responses back to the client.
    *   Dependencies: Application Layer.

2.  **Application Layer (`application`, `internal/*/application`)**:
    *   Orchestrates business logic and use cases.
    *   Contains application-specific services and use cases (e.g., `user_usecases.go`).
    *   Coordinates interactions between the Domain Layer and Infrastructure Layer.
    *   Dependencies: Domain Layer, Infrastructure Layer (via interfaces).

3.  **Domain Layer (`domain`, `internal/*/domain`)**:
    *   Contains the core business logic, entities, value objects, and interfaces for repositories and services.
    *   Independent of other layers.
    *   Dependencies: None (pure business logic).

4.  **Infrastructure Layer (`infrastructure`, `internal/*/infrastructure`)**:
    *   Provides implementations for interfaces defined in the Domain Layer.
    *   Handles external concerns like databases, external APIs, messaging queues, caching, etc.
    *   Examples: `infrastructure/database/postgres_user_repository.go`, `infrastructure/messaging`.
    *   Dependencies: Domain Layer (implements its interfaces).

5.  **Shared Layer (`shared`)**:
    *   Contains common utilities, cross-cutting concerns like authentication middleware, error handling, and constants.
    *   Dependencies: None (or very few, foundational).

## Dependency Flow

The dependencies generally flow inwards: Presentation -> Application -> Domain. The Infrastructure layer depends on the Domain layer by implementing its interfaces. The Shared layer provides utilities that can be used across multiple layers without introducing circular dependencies.

```
+-------------------+
|    Presentation   |
| (cmd/api, graph)  |
+---------+---------+
          |
          v
+---------+---------+
|    Application    |
| (application,     |
|  internal/*/app)  |
+---------+---------+
          |
          v
+---------+---------+
|      Domain       |
| (domain,          |
|  internal/*/domain)|
+---------+---------+
          ^
          |
+---------+---------+
|   Infrastructure  |
| (infrastructure,  |
|  internal/*/infra)|
+-------------------+

+-------------------+
|       Shared      |
|    (shared)       |
+-------------------+
  (Used by all layers)
```
