# Database Changes

This document describes the database schema changes and migrations implemented for the Chatear Backend project.

## Overview

The database integration uses Supabase (PostgreSQL) as the main database with a comprehensive schema that includes:

- User management with soft deletion
- Authentication tokens (refresh tokens, magic links)
- Action logging and audit trails
- User deletion management with capacity limits
- Email tracking and notifications

## Migration Files

### 0001_create_notifications_table.up.sql
- Creates the notifications table for system notifications
- Includes fields for type, recipient, subject, body, and timestamps

### 0002_create_refresh_tokens_table.up.sql
- Creates the refresh_tokens table for JWT refresh token management
- Includes token, expiration, revocation tracking, and user association

### 0003_create_main_schema.up.sql
- Creates the main application schema based on docs/database.md
- Includes tables: users, action_logs, deletion_capacity, email_sends, magic_links, user_deletion_cycles, user_deletions, user_logins

### 0004_create_indexes.up.sql
- Creates performance indexes for all tables
- Includes composite indexes for common query patterns
- Includes partial indexes for filtered queries

### 0005_create_functions_and_triggers.up.sql
- Creates database functions for business logic
- Creates triggers for automatic field updates
- Includes functions for preventing deleted user logins and updating timestamps

## Schema Details

### Users Table
- Primary user data storage with soft deletion support
- Email verification tracking
- Avatar URL support
- Deletion scheduling capabilities

### Action Logs Table
- Comprehensive audit trail for user actions
- JSON metadata support for flexible logging
- User association and timestamp tracking

### Magic Links Table
- Email verification and password reset tokens
- Automatic expiration and usage tracking
- Type-based categorization

### Refresh Tokens Table
- JWT refresh token management
- IP address and user agent tracking
- Revocation support

### User Deletion Management
- Scheduled user deletion with capacity limits
- Recovery token support for cancellation
- Cycle tracking for deletion attempts

## Database Adapter

The `infrastructure/db/postgres/adapter.go` provides:

- Connection pool management using pgx/v5
- Migration execution using golang-migrate
- Health checking capabilities
- Proper connection lifecycle management

## Repository Pattern

All database operations are abstracted through repository interfaces in the domain layer:

- `UserRepository` - User CRUD operations
- `ActionLogRepository` - Audit trail management
- `MagicLinkRepository` - Authentication token management
- `RefreshTokenRepository` - JWT token management
- `UserDeletionRepository` - Deletion scheduling
- `DeletionCapacityRepository` - Capacity management
- `EmailSendRepository` - Email tracking
- `UserLoginRepository` - Login attempt tracking

## Implementation Status

✅ Database adapter with migration support
✅ Complete schema migration files
✅ Domain entities for all database tables
✅ Repository interfaces in domain layer
✅ PostgreSQL repository implementations
✅ Indexes for performance optimization
✅ Database functions and triggers
✅ Soft deletion support
✅ Audit trail capabilities

## Usage

To run migrations:

```go
adapter, err := postgres.NewAdapter(databaseURL)
if err != nil {
    log.Fatal(err)
}
defer adapter.Close()

err = adapter.RunMigrations("./migrations/postgres")
if err != nil {
    log.Fatal(err)
}
```

## Security Considerations

- All user deletions are soft deletes to maintain data integrity
- Password hashes are never returned in API responses
- Audit trails are maintained for all user actions
- Token expiration is enforced at the database level
- Deleted users cannot log in (enforced by triggers)
