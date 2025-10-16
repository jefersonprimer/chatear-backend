package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserLogin represents a user login record
type UserLogin struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	IPAddress *string    `json:"ip_address,omitempty"`
	UserAgent *string    `json:"user_agent,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	Success   bool       `json:"success"`
}

// NewUserLogin creates a new user login record
func NewUserLogin(userID *uuid.UUID, ipAddress, userAgent *string, success bool) *UserLogin {
	return &UserLogin{
		ID:        uuid.New(),
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		Success:   success,
	}
}
