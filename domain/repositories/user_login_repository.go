package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// UserLoginRepository defines the interface for user login data operations
type UserLoginRepository interface {
	Create(ctx context.Context, userLogin *entities.UserLogin) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserLogin, error)
	GetByIPAddress(ctx context.Context, ipAddress string, limit, offset int) ([]*entities.UserLogin, error)
	GetRecentByUserID(ctx context.Context, userID uuid.UUID, since time.Time) ([]*entities.UserLogin, error)
	GetFailedLoginsByUserID(ctx context.Context, userID uuid.UUID, since time.Time) ([]*entities.UserLogin, error)
	CountFailedLoginsByUserID(ctx context.Context, userID uuid.UUID, since time.Time) (int, error)
	DeleteOldLogs(ctx context.Context, olderThan time.Time) error
}
