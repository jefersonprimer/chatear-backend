package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// ActionLogRepository defines the interface for action log data operations
type ActionLogRepository interface {
	Create(ctx context.Context, actionLog *entities.ActionLog) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.ActionLog, error)
	GetByAction(ctx context.Context, action string, limit, offset int) ([]*entities.ActionLog, error)
	GetByUserIDAndAction(ctx context.Context, userID uuid.UUID, action string, limit, offset int) ([]*entities.ActionLog, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entities.ActionLog, error)
	DeleteOldLogs(ctx context.Context, olderThan time.Time) error
}
