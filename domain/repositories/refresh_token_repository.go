package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// RefreshTokenRepository defines the interface for refresh token data operations
type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *entities.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.RefreshToken, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.RefreshToken, error)
	Update(ctx context.Context, refreshToken *entities.RefreshToken) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context, olderThan time.Time) error
	RevokeByUserID(ctx context.Context, userID uuid.UUID) error
	RevokeAllByUserID(userID uuid.UUID) error
}
