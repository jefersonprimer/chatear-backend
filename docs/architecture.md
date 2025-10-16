# Chatear Backend - Clean Architecture

## Overview

This project follows Clean Architecture principles with Domain-Driven Design (DDD) and SOLID principles to create a maintainable, testable, and scalable modular monolith.

## Architecture Layers

### 1. Domain Layer (`domain/`)
Contains the core business logic and entities. This layer is independent of any external concerns.

**Responsibilities:**
- Business entities and value objects
- Domain services
- Repository interfaces
- Domain events
- Business rules and validation

### 2. Application Layer (`application/`)
Contains use cases and application services that orchestrate the domain layer.

**Responsibilities:**
- Use case implementations
- Application services
- Command and query handlers
- DTOs and mappers
- Application events

### 3. Infrastructure Layer (`infrastructure/`)
Contains implementations of external concerns like databases, external APIs, and frameworks.

**Responsibilities:**
- Database implementations
- External service integrations
- Framework-specific code
- Configuration management
- Third-party library integrations

### 4. Presentation Layer (`presentation/`)
Contains the API endpoints, GraphQL resolvers, and web framework code.

**Responsibilities:**
- HTTP handlers
- GraphQL resolvers
- Request/response mapping
- Authentication middleware
- API documentation

### 5. Shared Layer (`shared/`)
Contains common utilities and cross-cutting concerns used across all layers.

**Responsibilities:**
- Common utilities
- Shared constants
- Cross-cutting concerns
- Common interfaces

## Module Organization

The project is organized into feature modules within the `internal/` directory:

```
internal/
├── user/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   └── presentation/
└── notification/
    ├── domain/
    ├── application/
    ├── infrastructure/
    └── presentation/
```

Each module is self-contained and follows the same Clean Architecture pattern.

## Technology Stack

- **Language:** Go 1.25.3
- **Web Framework:** Gin
- **API:** GraphQL (gqlgen)
- **Database:** PostgreSQL
- **Cache:** Redis
- **Message Queue:** NATS
- **Authentication:** JWT
- **Email:** SMTP

## Dependencies Flow

The dependency flow follows Clean Architecture principles:

```
Presentation → Application → Domain
     ↓              ↓
Infrastructure → Application → Domain
```

- Inner layers don't depend on outer layers
- Dependencies point inward
- Interfaces are defined in inner layers
- Implementations are in outer layers

## SOLID Principles

1. **Single Responsibility Principle (SRP):** Each class/module has one reason to change
2. **Open/Closed Principle (OCP):** Open for extension, closed for modification
3. **Liskov Substitution Principle (LSP):** Derived classes must be substitutable for their base classes
4. **Interface Segregation Principle (ISP):** Clients shouldn't depend on interfaces they don't use
5. **Dependency Inversion Principle (DIP):** Depend on abstractions, not concretions

## Getting Started

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Copy environment file: `cp env.example .env`
4. Configure your environment variables
5. Run the application: `go run cmd/api/main.go`