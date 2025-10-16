package infrastructure

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/infrastructure"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
)

// EmailRepository is a PostgreSQL implementation of the domain.EmailRepository.
type EmailRepository struct {
	DB *infrastructure.DB
}

// NewEmailRepository creates a new EmailRepository.
func NewEmailRepository(db *infrastructure.DB) *EmailRepository {
	return &EmailRepository{DB: db}
}

// CreateEmail creates a new email in the database.
func (r *EmailRepository) CreateEmail(ctx context.Context, email *domain.Email) error {
	query := `
		INSERT INTO email_sends (id, user_id, type)
		VALUES ($1, $2, $3)
		RETURNING sent_at`

	return r.DB.Pool.QueryRow(ctx, query,
		email.ID,
		email.UserID,
		email.Type,
	).Scan(&email.SentAt)
}

// GetEmailsByUserIDAndType retrieves all emails sent to a user of a specific type.
func (r *EmailRepository) GetEmailsByUserIDAndType(ctx context.Context, userID uuid.UUID, emailType string) ([]*domain.Email, error) {
	query := `
		SELECT id, user_id, type, sent_at
		FROM email_sends
		WHERE user_id = $1 AND type = $2 AND sent_at >= now() - interval '24 hours'`

	rows, err := r.DB.Pool.Query(ctx, query, userID, emailType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []*domain.Email
	for rows.Next() {
		email := &domain.Email{}
		if err := rows.Scan(&email.ID, &email.UserID, &email.Type, &email.SentAt); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, nil
}
