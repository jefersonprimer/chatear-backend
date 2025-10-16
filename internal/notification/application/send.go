package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jefersonprimer/chatear-backend/internal/notification/domain"
)

type EmailSender struct {
	repository domain.Repository
	sender     domain.Sender
}

func NewEmailSender(repository domain.Repository, sender domain.Sender) *EmailSender {
	return &EmailSender{
		repository: repository,
		sender:     sender,
	}
}

func (s *EmailSender) Send(ctx context.Context, recipient, subject, body, templateName string) (*domain.EmailSend, error) {
	emailSend := &domain.EmailSend{
		ID:           uuid.New().String(),
		Recipient:    recipient,
		Subject:      subject,
		Body:         body,
		TemplateName: templateName,
		SentAt:       time.Now(),
		Status:       "pending",
	}

	// Try to send the email
	if err := s.sender.Send(ctx, emailSend); err != nil {
		emailSend.Status = "failed"
		emailSend.ErrorMessage = err.Error()
		// Still save the failed attempt for logging purposes
		if saveErr := s.repository.Save(ctx, emailSend); saveErr != nil {
			return nil, saveErr
		}
		return nil, err
	}

	// Mark as sent and save
	emailSend.Status = "sent"
	if err := s.repository.Save(ctx, emailSend); err != nil {
		return nil, err
	}

	return emailSend, nil
}
