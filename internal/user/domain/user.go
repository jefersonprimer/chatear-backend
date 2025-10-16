package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID                uuid.UUID  `json:"id"`
	Name              string     `json:"name"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	IsEmailVerified   bool       `json:"is_email_verified"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	AvatarURL         *string    `json:"avatar_url,omitempty"`
	DeletionDueAt     *time.Time `json:"deletion_due_at,omitempty"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	IsDeleted         bool       `json:"is_deleted"`
}

// UserRepository defines the interface for interacting with user data.
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// BlacklistRepository defines the interface for managing blacklisted tokens.
type BlacklistRepository interface {
	Add(ctx context.Context, token string, expiration time.Duration) error
	Check(ctx context.Context, token string) (bool, error)
}
