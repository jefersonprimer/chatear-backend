package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Email represents an email sent to a user.
type Email struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	SentAt    time.Time `json:"sent_at"`
}

// EmailRepository defines the interface for interacting with email data.
type EmailRepository interface {
	CreateEmail(ctx context.Context, email *Email) error
	GetEmailsByUserIDAndType(ctx context.Context, userID uuid.UUID, emailType string) ([]*Email, error)
}
