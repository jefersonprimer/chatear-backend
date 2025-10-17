package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"github.com/jefersonprimer/chatear-backend/internal/user/infrastructure"
	"github.com/jefersonprimer/chatear-backend/shared/events"
	"github.com/jefersonprimer/chatear-backend/shared/util"
)

// ResendVerificationEmail is a use case for resending the verification email.
type ResendVerificationEmail struct {
	UserRepository  domain.UserRepository
	EmailRepository domain.EmailRepository
	TokenRepository infrastructure.TokenRepository
	EventBus        domain.EventBus
	AppURL          string
	MaxEmailsPerDay int
}

// NewResendVerificationEmail creates a new ResendVerificationEmail use case.
func NewResendVerificationEmail(userRepository domain.UserRepository, emailRepository domain.EmailRepository, tokenRepository infrastructure.TokenRepository, eventBus domain.EventBus, appURL string, maxEmailsPerDay int) *ResendVerificationEmail {
	if maxEmailsPerDay == 0 {
		maxEmailsPerDay = defaultMaxEmailsPerDay
	}
	return &ResendVerificationEmail{
		UserRepository:  userRepository,
		EmailRepository: emailRepository,
		TokenRepository: tokenRepository,
		EventBus:        eventBus,
		AppURL:          appURL,
		MaxEmailsPerDay: maxEmailsPerDay,
	}
}

// Execute finds a user by email, generates a new verification token, and sends it.
func (uc *ResendVerificationEmail) Execute(ctx context.Context, email string) error {
	user, err := uc.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsEmailVerified {
		return errors.New("email already verified")
	}

	emails, err := uc.EmailRepository.GetEmailsByUserIDAndType(ctx, user.ID, "verification")
	if err != nil {
		return err
	}

	if len(emails) >= uc.MaxEmailsPerDay {
		return errors.New("email limit exceeded")
	}

	token, err := util.GenerateRandomToken()
	if err != nil {
		return err
	}

	if err := uc.TokenRepository.Set(ctx, fmt.Sprintf("verification:%s", token), user.ID.String(), 15*time.Minute); err != nil {
		return err
	}

	emailRequest := events.EmailSendRequest{
		Recipient: user.Email,
		Subject:   "Email Verification",
		Body:      fmt.Sprintf("Click here to verify your email: %s/verify-email?token=%s", uc.AppURL, token),
	}
	emailDataBytes, err := json.Marshal(emailRequest)
	if err != nil {
		return err
	}

	if err := uc.EventBus.Publish(ctx, &domain.Event{Subject: "email.send", Data: emailDataBytes}); err != nil {
		return err
	}

	if err := uc.EmailRepository.CreateEmail(ctx, &domain.Email{ID: uuid.New(), UserID: user.ID, Type: "verification"}); err != nil {
		return err
	}

	return nil
}
