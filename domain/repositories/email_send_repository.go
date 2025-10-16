package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// EmailSendRepository defines the interface for email send data operations
type EmailSendRepository interface {
	Create(ctx context.Context, emailSend *entities.EmailSend) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.EmailSend, error)
	GetByType(ctx context.Context, emailType entities.EmailType, limit, offset int) ([]*entities.EmailSend, error)
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, emailType entities.EmailType, limit, offset int) ([]*entities.EmailSend, error)
	GetRecentByUserIDAndType(ctx context.Context, userID uuid.UUID, emailType entities.EmailType, since time.Time) ([]*entities.EmailSend, error)
	CountByUserIDAndType(ctx context.Context, userID uuid.UUID, emailType entities.EmailType, since time.Time) (int, error)
}
