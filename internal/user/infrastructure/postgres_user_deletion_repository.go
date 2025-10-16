package infrastructure

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/infrastructure"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
)

// UserDeletionRepository is a PostgreSQL implementation of the domain.UserDeletionRepository.
type UserDeletionRepository struct {
	DB *infrastructure.DB
}

// NewUserDeletionRepository creates a new UserDeletionRepository.
func NewUserDeletionRepository(db *infrastructure.DB) *UserDeletionRepository {
	return &UserDeletionRepository{DB: db}
}

// CreateUserDeletion creates a new user deletion in the database.
func (r *UserDeletionRepository) CreateUserDeletion(ctx context.Context, userDeletion *domain.UserDeletion) error {
	query := `
		INSERT INTO user_deletions (id, user_id, scheduled_date, status)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`

	return r.DB.Pool.QueryRow(ctx, query,
		userDeletion.ID,
		userDeletion.UserID,
		userDeletion.ScheduledDate,
		userDeletion.Status,
	).Scan(&userDeletion.CreatedAt)
}

// GetUserDeletionsByDate retrieves all user deletions for a specific date.
func (r *UserDeletionRepository) GetUserDeletionsByDate(ctx context.Context, date time.Time) ([]*domain.UserDeletion, error) {
	query := `
		SELECT id, user_id, scheduled_date, executed, created_at, status, token, token_expires_at, recovery_token, recovery_token_expires_at
		FROM user_deletions
		WHERE scheduled_date = $1`

	rows, err := r.DB.Pool.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userDeletions []*domain.UserDeletion
	for rows.Next() {
		userDeletion := &domain.UserDeletion{}
		if err := rows.Scan(
			&userDeletion.ID,
			&userDeletion.UserID,
			&userDeletion.ScheduledDate,
			&userDeletion.Executed,
			&userDeletion.CreatedAt,
			&userDeletion.Status,
			&userDeletion.Token,
			&userDeletion.TokenExpiresAt,
			&userDeletion.RecoveryToken,
			&userDeletion.RecoveryTokenExpiresAt,
		); err != nil {
			return nil, err
		}
		userDeletions = append(userDeletions, userDeletion)
	}

	return userDeletions, nil
}