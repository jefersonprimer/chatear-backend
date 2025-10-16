package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

// RefreshToken represents a refresh token in the system.
type RefreshToken struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	Revoked   bool       `json:"revoked"`
}

// RefreshTokenRepository defines the interface for managing refresh tokens.
type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshTokenByToken(ctx context.Context, token string) (*RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, token *RefreshToken) error
	RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
}