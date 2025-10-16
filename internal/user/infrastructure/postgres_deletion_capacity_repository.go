package infrastructure

import (
	"context"
	"time"

	"github.com/jefersonprimer/chatear-backend/infrastructure"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
)

// DeletionCapacityRepository is a PostgreSQL implementation of the domain.DeletionCapacityRepository.
type DeletionCapacityRepository struct {
	DB *infrastructure.DB
}

// NewDeletionCapacityRepository creates a new DeletionCapacityRepository.
func NewDeletionCapacityRepository(db *infrastructure.DB) *DeletionCapacityRepository {
	return &DeletionCapacityRepository{DB: db}
}

// GetDeletionCapacity retrieves the deletion capacity for a given day.
func (r *DeletionCapacityRepository) GetDeletionCapacity(ctx context.Context, day time.Time) (*domain.DeletionCapacity, error) {
	capacity := &domain.DeletionCapacity{}
	query := `
		SELECT day, count, max_limit, updated_at
		FROM deletion_capacity
		WHERE day = $1`

	err := r.DB.Pool.QueryRow(ctx, query, day).Scan(
		&capacity.Day,
		&capacity.Count,
		&capacity.MaxLimit,
		&capacity.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return capacity, nil
}

// IncrementDeletionCapacity increments the deletion capacity for a given day.
func (r *DeletionCapacityRepository) IncrementDeletionCapacity(ctx context.Context, day time.Time) error {
	query := `
		INSERT INTO deletion_capacity (day, count)
		VALUES ($1, 1)
		ON CONFLICT (day) DO UPDATE
		SET count = deletion_capacity.count + 1`

	_, err := r.DB.Pool.Exec(ctx, query, day)
	return err
}
