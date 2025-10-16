package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// UserDeletionRepository defines the interface for user deletion data operations
type UserDeletionRepository interface {
	Create(ctx context.Context, userDeletion *entities.UserDeletion) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.UserDeletion, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserDeletion, error)
	GetByToken(ctx context.Context, token string) (*entities.UserDeletion, error)
	GetByRecoveryToken(ctx context.Context, recoveryToken string) (*entities.UserDeletion, error)
	GetScheduledDeletions(ctx context.Context, scheduledDate time.Time) ([]*entities.UserDeletion, error)
	GetByStatus(ctx context.Context, status entities.UserDeletionStatus, limit, offset int) ([]*entities.UserDeletion, error)
	Update(ctx context.Context, userDeletion *entities.UserDeletion) error
	Delete(ctx context.Context, id uuid.UUID) error
	CancelByUserID(ctx context.Context, userID uuid.UUID) error
}
