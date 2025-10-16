package repositories

import (
	"context"
	"time"

	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// DeletionCapacityRepository defines the interface for deletion capacity data operations
type DeletionCapacityRepository interface {
	GetByDate(ctx context.Context, date time.Time) (*entities.DeletionCapacity, error)
	Create(ctx context.Context, capacity *entities.DeletionCapacity) error
	Update(ctx context.Context, capacity *entities.DeletionCapacity) error
	IncrementCount(ctx context.Context, date time.Time) error
	GetAvailableCapacity(ctx context.Context, date time.Time) (int, error)
}
