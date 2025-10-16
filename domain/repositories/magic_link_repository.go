package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// MagicLinkRepository defines the interface for magic link data operations
type MagicLinkRepository interface {
	Create(ctx context.Context, magicLink *entities.MagicLink) error
	GetByToken(ctx context.Context, token string) (*entities.MagicLink, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.MagicLink, error)
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, linkType entities.MagicLinkType, limit, offset int) ([]*entities.MagicLink, error)
	GetActiveByUserIDAndType(ctx context.Context, userID uuid.UUID, linkType entities.MagicLinkType) ([]*entities.MagicLink, error)
	Update(ctx context.Context, magicLink *entities.MagicLink) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context, olderThan time.Time) error
	RevokeByUserIDAndType(ctx context.Context, userID uuid.UUID, linkType entities.MagicLinkType) error
}
