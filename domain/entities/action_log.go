package entities

import (
	"time"

	"github.com/google/uuid"
)

// ActionLog represents an action log entry in the system
type ActionLog struct {
	ID        uuid.UUID       `json:"id"`
	UserID    *uuid.UUID      `json:"user_id,omitempty"`
	Action    string          `json:"action"`
	CreatedAt time.Time       `json:"created_at"`
	Meta      map[string]any  `json:"meta,omitempty"`
}

// NewActionLog creates a new action log entry
func NewActionLog(userID *uuid.UUID, action string, meta map[string]any) *ActionLog {
	return &ActionLog{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		CreatedAt: time.Now(),
		Meta:      meta,
	}
}
