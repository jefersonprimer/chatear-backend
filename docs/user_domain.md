# User Domain

This document describes the structure and flows of the user domain, which implements a comprehensive user management system with authentication, authorization, and account lifecycle management.

## Folder Structure

```
internal/user/
├── application/           # Use cases (application layer)
│   ├── delete.go         # Account deletion use case
│   ├── login.go          # User login use case
│   ├── logout.go         # User logout use case
│   ├── password_recovery.go  # Password recovery use case
│   ├── register.go       # User registration use case
│   └── verify_token.go   # Token verification use case
├── domain/               # Domain entities and interfaces
│   ├── deletion_capacity.go  # Deletion capacity entity
│   ├── email.go          # Email tracking entity
│   ├── event_bus.go      # Event bus interface
│   ├── refresh_token.go  # Refresh token entity
│   ├── token_service.go  # Token service interface
│   ├── user.go           # User entity and repository interface
│   └── user_deletion.go  # User deletion entity
├── infrastructure/       # Infrastructure implementations
│   ├── jwt_service.go    # JWT token service implementation
│   ├── nats_event_bus.go # NATS event bus implementation
│   ├── postgres_deletion_capacity_repository.go
│   ├── postgres_email_repository.go
│   ├── postgres_refresh_token_repository.go
│   ├── postgres_repository.go
│   ├── postgres_user_deletion_repository.go
│   └── redis_cache.go    # Redis token cache implementation
└── presentation/         # Presentation layer
    ├── gin_handlers.go   # REST API handlers
    ├── graphql_resolvers.go  # GraphQL resolvers
    └── schema.graphqls   # GraphQL schema
```

## Architecture

The user domain follows Clean Architecture principles with clear separation of concerns:

### Domain Layer
- **Entities**: Core business objects (User, RefreshToken, Email, UserDeletion, DeletionCapacity)
- **Interfaces**: Repository and service contracts
- **Value Objects**: Business rules and constraints
- **Events**: Domain events for decoupled communication

### Application Layer
- **Use Cases**: Business logic orchestration
  - `RegisterUser`: User registration with email verification
  - `LoginUser`: Authentication with JWT and refresh tokens
  - `LogoutUser`: Token invalidation
  - `PasswordRecovery`: Password reset via magic links
  - `DeleteUser`: Account deletion scheduling
  - `VerifyToken`: Token verification for various purposes

### Infrastructure Layer
- **Repositories**: Data persistence implementations
  - PostgreSQL for user data, refresh tokens, email tracking
  - Redis for temporary tokens (verification, password reset)
- **Services**: External service integrations
  - JWT service for token generation/validation
  - NATS for event publishing
  - Redis for caching

### Presentation Layer
- **REST API**: Gin handlers for HTTP endpoints
- **GraphQL API**: Resolvers for GraphQL queries/mutations
- **Schema**: GraphQL type definitions

## Use Cases and Flows

### 1. User Registration
**Endpoint**: `POST /register` or GraphQL `register` mutation

**Flow**:
1. User provides name, email, and password
2. Password is hashed using bcrypt
3. User record is created in PostgreSQL
4. Email limit check (max 2 verification emails per day)
5. Verification token generated and stored in Redis (15-minute TTL)
6. `email.send` event published to NATS
7. Email record created for tracking

**Business Rules**:
- Email must be unique
- Password must be hashed before storage
- Maximum 2 verification emails per user per day
- Verification tokens expire in 15 minutes

### 2. User Login
**Endpoint**: `POST /login` or GraphQL `login` mutation

**Flow**:
1. User provides email and password
2. User retrieved from database by email
3. Password verified using bcrypt
4. JWT access token generated (15-minute expiry)
5. Refresh token generated and stored in PostgreSQL (7-day expiry)
6. Returns both tokens to client

**Response**:
```json
{
  "accessToken": "jwt_token_here",
  "refreshToken": "refresh_token_here"
}
```

### 3. User Logout
**Endpoint**: `POST /logout` or GraphQL `logout` mutation

**Flow**:
1. User provides refresh token
2. Refresh token retrieved from database
3. Token marked as revoked with timestamp
4. Token updated in database

### 4. Password Recovery
**Endpoint**: `POST /password-recovery` or GraphQL `passwordRecovery` mutation

**Flow**:
1. User provides email address
2. User retrieved from database
3. Password reset token generated
4. Token stored in Redis (15-minute TTL)
5. `email.send` event published to NATS
6. Magic link sent to user's email

### 5. Account Deletion
**Endpoint**: `DELETE /users/:id` or GraphQL `deleteUser` mutation

**Flow**:
1. User ID provided for deletion
2. Daily deletion limit checked (max 10 per day)
3. Deletion scheduled for 90 days in the future
4. `user.delete` event published to NATS
5. Daily deletion count incremented

**Business Rules**:
- Maximum 10 account deletions per day
- Deletion scheduled with 90-day delay
- User data marked for deletion, not immediately removed

### 6. Token Verification
**Endpoint**: `GET /verify-email?token=...` or `GET /reset-password?token=...`

**Flow**:
1. Token retrieved from Redis
2. User ID extracted from token
3. User retrieved from database
4. Action performed based on token type:
   - `verification`: Mark email as verified
   - `password-reset`: Allow password reset
5. Token removed from Redis

## Data Storage

### PostgreSQL Tables
- `users`: User account information
- `refresh_tokens`: JWT refresh tokens
- `email_sends`: Email tracking for rate limiting
- `user_deletions`: Scheduled account deletions
- `deletion_capacity`: Daily deletion limits

### Redis Cache
- `verification:{token}`: Email verification tokens (15min TTL)
- `password-reset:{token}`: Password reset tokens (15min TTL)

## Events

### NATS Events
- `email.send`: Triggered when email needs to be sent
- `user.delete`: Triggered when user deletion is scheduled

## API Endpoints

### REST API (Gin)
- `POST /register` - User registration
- `POST /login` - User login
- `POST /logout` - User logout
- `POST /password-recovery` - Password recovery
- `DELETE /users/:id` - Account deletion
- `GET /verify-email?token=...` - Email verification
- `GET /reset-password?token=...` - Password reset

### GraphQL API
```graphql
type User {
    id: ID!
    name: String!
    email: String!
    isEmailVerified: Boolean!
    avatarUrl: String
    createdAt: Time!
    updatedAt: Time!
    lastLoginAt: Time
    isDeleted: Boolean!
}

type LoginResponse {
    accessToken: String!
    refreshToken: String!
}

type Query {
    me: User
}

type Mutation {
    register(input: RegisterUserInput!): User
    login(input: LoginInput!): LoginResponse
    logout(input: LogoutInput!): Boolean
    deleteUser(id: ID!): Boolean
}
```

## Security Features

- Password hashing with bcrypt
- JWT access tokens with 15-minute expiry
- Refresh tokens with 7-day expiry
- Rate limiting on email sending (2 per day)
- Rate limiting on account deletions (10 per day)
- Token-based email verification
- Magic link password recovery
- Scheduled account deletion with recovery period

## Dependencies

- **PostgreSQL**: Primary database for user data
- **Redis**: Token cache and temporary storage
- **NATS**: Event messaging system
- **JWT**: Token-based authentication
- **bcrypt**: Password hashing
- **Gin**: HTTP web framework
- **GraphQL**: API query language
