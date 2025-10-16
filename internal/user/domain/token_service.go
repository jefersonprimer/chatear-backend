package domain

import (
	"context"

	"github.com/google/uuid"
)

// TokenService defines the interface for creating and validating tokens.
type TokenService interface {
	CreateAccessToken(ctx context.Context, user *User) (string, error)
	CreateRefreshToken(ctx context.Context, user *User) (string, error)
	VerifyToken(ctx context.Context, tokenString string) (uuid.UUID, error)
}
