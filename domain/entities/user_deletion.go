package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserDeletionStatus represents the status of a user deletion
type UserDeletionStatus string

const (
	UserDeletionStatusQueued    UserDeletionStatus = "queued"
	UserDeletionStatusScheduled UserDeletionStatus = "scheduled"
	UserDeletionStatusExecuted  UserDeletionStatus = "executed"
	UserDeletionStatusCancelled UserDeletionStatus = "cancelled"
)

// UserDeletion represents a user deletion request
type UserDeletion struct {
	ID                      uuid.UUID           `json:"id"`
	UserID                  uuid.UUID           `json:"user_id"`
	ScheduledDate           time.Time           `json:"scheduled_date"`
	Executed                bool                `json:"executed"`
	CreatedAt               time.Time           `json:"created_at"`
	Status                  UserDeletionStatus  `json:"status"`
	Token                   *string             `json:"token,omitempty"`
	TokenExpiresAt          *time.Time          `json:"token_expires_at,omitempty"`
	RecoveryToken           *string             `json:"recovery_token,omitempty"`
	RecoveryTokenExpiresAt  *time.Time          `json:"recovery_token_expires_at,omitempty"`
}

// NewUserDeletion creates a new user deletion request
func NewUserDeletion(userID uuid.UUID, scheduledDate time.Time) *UserDeletion {
	return &UserDeletion{
		ID:            uuid.New(),
		UserID:        userID,
		ScheduledDate: scheduledDate,
		Executed:      false,
		CreatedAt:     time.Now(),
		Status:        UserDeletionStatusQueued,
	}
}

// MarkAsScheduled marks the deletion as scheduled
func (ud *UserDeletion) MarkAsScheduled() {
	ud.Status = UserDeletionStatusScheduled
}

// MarkAsExecuted marks the deletion as executed
func (ud *UserDeletion) MarkAsExecuted() {
	ud.Status = UserDeletionStatusExecuted
	ud.Executed = true
}

// MarkAsCancelled marks the deletion as cancelled
func (ud *UserDeletion) MarkAsCancelled() {
	ud.Status = UserDeletionStatusCancelled
}
