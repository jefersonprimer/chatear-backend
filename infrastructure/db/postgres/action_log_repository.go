package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
	"github.com/jefersonprimer/chatear-backend/domain/repositories"
)

// actionLogRepository implements the ActionLogRepository interface
type actionLogRepository struct {
	adapter *Adapter
}

// NewActionLogRepository creates a new action log repository
func NewActionLogRepository(adapter *Adapter) repositories.ActionLogRepository {
	return &actionLogRepository{
		adapter: adapter,
	}
}

// Create creates a new action log entry
func (r *actionLogRepository) Create(ctx context.Context, actionLog *entities.ActionLog) error {
	query := `
		INSERT INTO action_logs (id, user_id, action, created_at, meta)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.adapter.Pool.Exec(ctx, query,
		actionLog.ID,
		actionLog.UserID,
		actionLog.Action,
		actionLog.CreatedAt,
		actionLog.Meta,
	)

	return err
}

// GetByUserID retrieves action logs by user ID
func (r *actionLogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.ActionLog, error) {
	query := `
		SELECT id, user_id, action, created_at, meta
		FROM action_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.adapter.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActionLogs(rows)
}

// GetByAction retrieves action logs by action type
func (r *actionLogRepository) GetByAction(ctx context.Context, action string, limit, offset int) ([]*entities.ActionLog, error) {
	query := `
		SELECT id, user_id, action, created_at, meta
		FROM action_logs
		WHERE action = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.adapter.Pool.Query(ctx, query, action, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActionLogs(rows)
}

// GetByUserIDAndAction retrieves action logs by user ID and action type
func (r *actionLogRepository) GetByUserIDAndAction(ctx context.Context, userID uuid.UUID, action string, limit, offset int) ([]*entities.ActionLog, error) {
	query := `
		SELECT id, user_id, action, created_at, meta
		FROM action_logs
		WHERE user_id = $1 AND action = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.adapter.Pool.Query(ctx, query, userID, action, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActionLogs(rows)
}

// GetByDateRange retrieves action logs within a date range
func (r *actionLogRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entities.ActionLog, error) {
	query := `
		SELECT id, user_id, action, created_at, meta
		FROM action_logs
		WHERE created_at >= $1 AND created_at <= $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.adapter.Pool.Query(ctx, query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActionLogs(rows)
}

// DeleteOldLogs deletes action logs older than the specified time
func (r *actionLogRepository) DeleteOldLogs(ctx context.Context, olderThan time.Time) error {
	query := `DELETE FROM action_logs WHERE created_at < $1`
	_, err := r.adapter.Pool.Exec(ctx, query, olderThan)
	return err
}

// scanActionLogs scans rows into action log entities
func (r *actionLogRepository) scanActionLogs(rows interface{}) ([]*entities.ActionLog, error) {
	// This would need to be implemented based on the actual pgx rows interface
	// For now, returning empty slice as placeholder
	return []*entities.ActionLog{}, nil
}
