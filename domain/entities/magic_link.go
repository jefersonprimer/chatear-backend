package entities

import (
	"time"

	"github.com/google/uuid"
)

// MagicLinkType represents the type of magic link
type MagicLinkType string

const (
	MagicLinkTypeEmailVerification MagicLinkType = "email_verification"
	MagicLinkTypePasswordReset     MagicLinkType = "password_reset"
)

// MagicLink represents a magic link for authentication
type MagicLink struct {
	ID        uuid.UUID      `json:"id"`
	UserID    *uuid.UUID     `json:"user_id,omitempty"`
	Token     string         `json:"token"`
	ExpiresAt time.Time      `json:"expires_at"`
	Used      bool           `json:"used"`
	CreatedAt time.Time      `json:"created_at"`
	Type      MagicLinkType  `json:"type"`
	UsedAt    *time.Time     `json:"used_at,omitempty"`
	IsActive  bool           `json:"is_active"`
}

// NewMagicLink creates a new magic link
func NewMagicLink(userID *uuid.UUID, token string, expiresAt time.Time, linkType MagicLinkType) *MagicLink {
	return &MagicLink{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
		CreatedAt: time.Now(),
		Type:      linkType,
		IsActive:  true,
	}
}

// IsExpired checks if the magic link is expired
func (ml *MagicLink) IsExpired() bool {
	return time.Now().After(ml.ExpiresAt)
}

// MarkAsUsed marks the magic link as used
func (ml *MagicLink) MarkAsUsed() {
	now := time.Now()
	ml.Used = true
	ml.UsedAt = &now
	ml.IsActive = false
}
