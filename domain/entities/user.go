package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity in the domain
type User struct {
	ID                uuid.UUID  `json:"id"`
	Name              string     `json:"name"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	IsEmailVerified   bool       `json:"is_email_verified"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	AvatarURL         *string    `json:"avatar_url,omitempty"`
	DeletionDueAt     *time.Time `json:"deletion_due_at,omitempty"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	IsDeleted         bool       `json:"is_deleted"`
}

// NewUser creates a new user entity
func NewUser(name, email, passwordHash string) *User {
	return &User{
		ID:              uuid.New(),
		Name:            name,
		Email:           email,
		PasswordHash:    passwordHash,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsEmailVerified: false,
		IsDeleted:       false,
	}
}

// Validate validates the user entity
func (u *User) Validate() error {
	// Add validation logic here
	return nil
}

// MarkAsDeleted marks the user as deleted
func (u *User) MarkAsDeleted() {
	now := time.Now()
	u.IsDeleted = true
	u.DeletedAt = &now
	u.UpdatedAt = now
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// VerifyEmail marks the user's email as verified
func (u *User) VerifyEmail() {
	u.IsEmailVerified = true
	u.UpdatedAt = time.Now()
}