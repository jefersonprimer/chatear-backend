package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserDeletionCycle represents a user's deletion cycle tracking
type UserDeletionCycle struct {
	UserID       uuid.UUID  `json:"user_id"`
	Cycles       int        `json:"cycles"`
	LastCycleAt  *time.Time `json:"last_cycle_at,omitempty"`
}

// NewUserDeletionCycle creates a new user deletion cycle
func NewUserDeletionCycle(userID uuid.UUID) *UserDeletionCycle {
	return &UserDeletionCycle{
		UserID: userID,
		Cycles: 0,
	}
}

// IncrementCycle increments the deletion cycle count
func (udc *UserDeletionCycle) IncrementCycle() {
	now := time.Now()
	udc.Cycles++
	udc.LastCycleAt = &now
}
