package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// UserDeletionCycleRepository defines the interface for user deletion cycle data operations
type UserDeletionCycleRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserDeletionCycle, error)
	Create(ctx context.Context, cycle *entities.UserDeletionCycle) error
	Update(ctx context.Context, cycle *entities.UserDeletionCycle) error
	IncrementCycle(ctx context.Context, userID uuid.UUID) error
}
