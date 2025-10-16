package application

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

// EmailSendRequest represents a request to send an email
type EmailSendRequest struct {
	Recipient    string `json:"recipient"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	TemplateName string `json:"template_name,omitempty"`
}

// EmailService provides a high-level interface for sending emails
type EmailService struct {
	natsConn *nats.Conn
}

// NewEmailService creates a new EmailService instance
func NewEmailService(natsConn *nats.Conn) *EmailService {
	return &EmailService{
		natsConn: natsConn,
	}
}

// SendWelcomeEmail sends a welcome email to a user
func (s *EmailService) SendWelcomeEmail(ctx context.Context, recipient, userName string) error {
	subject := "Welcome to Chatear!"
	body := "Welcome to Chatear! We're excited to have you on board."
	
	return s.publishEmailEvent(ctx, recipient, subject, body, "welcome")
}

// SendMagicLinkEmail sends a magic link email for authentication
func (s *EmailService) SendMagicLinkEmail(ctx context.Context, recipient, magicLink string) error {
	subject := "Your Magic Link"
	body := "Click the following link to sign in: " + magicLink
	
	return s.publishEmailEvent(ctx, recipient, subject, body, "magic_link")
}

// SendCustomEmail sends a custom email with optional template
func (s *EmailService) SendCustomEmail(ctx context.Context, recipient, subject, body, templateName string) error {
	return s.publishEmailEvent(ctx, recipient, subject, body, templateName)
}

// publishEmailEvent publishes an email send event to NATS
func (s *EmailService) publishEmailEvent(ctx context.Context, recipient, subject, body, templateName string) error {
	request := EmailSendRequest{
		Recipient:    recipient,
		Subject:      subject,
		Body:         body,
		TemplateName: templateName,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return s.natsConn.Publish("email.send", data)
}