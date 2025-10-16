package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetDeletedUsers(ctx context.Context, limit, offset int) ([]*entities.User, error)
	GetByEmailVerified(ctx context.Context, verified bool, limit, offset int) ([]*entities.User, error)
	SearchByName(ctx context.Context, name string, limit, offset int) ([]*entities.User, error)
}