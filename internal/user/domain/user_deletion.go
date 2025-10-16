package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserDeletion represents a user deletion request.
type UserDeletion struct {
	ID                     uuid.UUID  `json:"id"`
	UserID                 uuid.UUID  `json:"user_id"`
	ScheduledDate          time.Time  `json:"scheduled_date"`
	Executed               bool       `json:"executed"`
	CreatedAt              time.Time  `json:"created_at"`
	Status                 string     `json:"status"`
	Token                  *string    `json:"token,omitempty"`
	TokenExpiresAt         *time.Time `json:"token_expires_at,omitempty"`
	RecoveryToken          *string    `json:"recovery_token,omitempty"`
	RecoveryTokenExpiresAt *time.Time `json:"recovery_token_expires_at,omitempty"`
}

// UserDeletionRepository defines the interface for interacting with user deletion data.
type UserDeletionRepository interface {
	CreateUserDeletion(ctx context.Context, userDeletion *UserDeletion) error
	GetUserDeletionsByDate(ctx context.Context, date time.Time) ([]*UserDeletion, error)
}
