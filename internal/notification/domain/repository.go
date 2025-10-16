package domain

import "context"

type Repository interface {
	Save(ctx context.Context, emailSend *EmailSend) error
	GetByID(ctx context.Context, id string) (*EmailSend, error)
	GetByRecipient(ctx context.Context, recipient string, limit int) ([]*EmailSend, error)
}
