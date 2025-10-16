package entities

import (
	"time"

	"github.com/google/uuid"
)

// EmailType represents the type of email being sent
type EmailType string

const (
	EmailTypeVerification  EmailType = "verification"
	EmailTypePasswordReset EmailType = "password_reset"
)

// EmailSend represents an email send record
type EmailSend struct {
	ID      uuid.UUID `json:"id"`
	UserID  *uuid.UUID `json:"user_id,omitempty"`
	Type    EmailType `json:"type"`
	SentAt  time.Time `json:"sent_at"`
}

// NewEmailSend creates a new email send record
func NewEmailSend(userID *uuid.UUID, emailType EmailType) *EmailSend {
	return &EmailSend{
		ID:     uuid.New(),
		UserID: userID,
		Type:   emailType,
		SentAt: time.Now(),
	}
}
