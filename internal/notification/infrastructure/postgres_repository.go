package infrastructure

import (
	"context"

	"github.com/jefersonprimer/chatear-backend/infrastructure"
	"github.com/jefersonprimer/chatear-backend/internal/notification/domain"
)

type PostgresRepository struct {
	DB *infrastructure.DB
}

func NewPostgresRepository(db *infrastructure.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) Save(ctx context.Context, emailSend *domain.EmailSend) error {
	query := `
		INSERT INTO email_sends (id, recipient, subject, body, template_name, sent_at, error_message, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at`

	return r.DB.Pool.QueryRow(ctx, query,
		emailSend.ID,
		emailSend.Recipient,
		emailSend.Subject,
		emailSend.Body,
		emailSend.TemplateName,
		emailSend.SentAt,
		emailSend.ErrorMessage,
		emailSend.Status,
	).Scan(&emailSend.CreatedAt)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.EmailSend, error) {
	query := `
		SELECT id, recipient, subject, body, template_name, sent_at, created_at, error_message, status
		FROM email_sends
		WHERE id = $1`

	emailSend := &domain.EmailSend{}
	err := r.DB.Pool.QueryRow(ctx, query, id).Scan(
		&emailSend.ID,
		&emailSend.Recipient,
		&emailSend.Subject,
		&emailSend.Body,
		&emailSend.TemplateName,
		&emailSend.SentAt,
		&emailSend.CreatedAt,
		&emailSend.ErrorMessage,
		&emailSend.Status,
	)

	if err != nil {
		return nil, err
	}

	return emailSend, nil
}

func (r *PostgresRepository) GetByRecipient(ctx context.Context, recipient string, limit int) ([]*domain.EmailSend, error) {
	query := `
		SELECT id, recipient, subject, body, template_name, sent_at, created_at, error_message, status
		FROM email_sends
		WHERE recipient = $1
		ORDER BY sent_at DESC
		LIMIT $2`

	rows, err := r.DB.Pool.Query(ctx, query, recipient, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emailSends []*domain.EmailSend
	for rows.Next() {
		emailSend := &domain.EmailSend{}
		err := rows.Scan(
			&emailSend.ID,
			&emailSend.Recipient,
			&emailSend.Subject,
			&emailSend.Body,
			&emailSend.TemplateName,
			&emailSend.SentAt,
			&emailSend.CreatedAt,
			&emailSend.ErrorMessage,
			&emailSend.Status,
		)
		if err != nil {
			return nil, err
		}
		emailSends = append(emailSends, emailSend)
	}

	return emailSends, nil
}
