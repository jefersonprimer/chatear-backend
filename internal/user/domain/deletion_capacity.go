package domain

import (
	"context"
	"time"
)

// DeletionCapacity represents the deletion capacity for a given day.
type DeletionCapacity struct {
	Day       time.Time `json:"day"`
	Count     int       `json:"count"`
	MaxLimit  int       `json:"max_limit"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeletionCapacityRepository defines the interface for interacting with deletion capacity data.
type DeletionCapacityRepository interface {
	GetDeletionCapacity(ctx context.Context, day time.Time) (*DeletionCapacity, error)
	IncrementDeletionCapacity(ctx context.Context, day time.Time) error
}
