package events

import "time"

// UserRegisteredEvent is published when a new user registers.
type UserRegisteredEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}
