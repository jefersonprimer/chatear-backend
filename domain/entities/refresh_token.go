package entities

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token for authentication
type RefreshToken struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	Revoked   bool       `json:"revoked"`
	IPAddress *string    `json:"ip_address,omitempty"`
	UserAgent *string    `json:"user_agent,omitempty"`
}

// NewRefreshToken creates a new refresh token
func NewRefreshToken(userID *uuid.UUID, token string, expiresAt time.Time, ipAddress, userAgent *string) *RefreshToken {
	return &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Revoked:   false,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
}

// IsExpired checks if the refresh token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// Revoke marks the refresh token as revoked
func (rt *RefreshToken) Revoke() {
	rt.Revoked = true
}
