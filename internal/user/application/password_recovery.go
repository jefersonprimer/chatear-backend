import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"github.com/jefersonprimer/chatear-backend/internal/user/infrastructure"
	"github.com/jefersonprimer/chatear-backend/shared/events"
	"github.com/jefersonprimer/chatear-backend/shared/util"
)

// PasswordRecovery is a use case for recovering a user's password.
type PasswordRecovery struct {
	UserRepository  domain.UserRepository
	TokenRepository infrastructure.TokenRepository
	EventBus        domain.EventBus
}

// NewPasswordRecovery creates a new PasswordRecovery use case.
func NewPasswordRecovery(userRepository domain.UserRepository, tokenRepository infrastructure.TokenRepository, eventBus domain.EventBus) *PasswordRecovery {
	return &PasswordRecovery{
		UserRepository:  userRepository,
		TokenRepository: tokenRepository,
		EventBus:        eventBus,
	}
}

// Execute sends a password recovery email to the user.
func (uc *PasswordRecovery) Execute(ctx context.Context, email string) error {
	user, err := uc.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	token, err := util.GenerateRandomToken()
	if err != nil {
		return err
	}

	if err := uc.TokenRepository.Set(ctx, fmt.Sprintf("password-reset:%s", token), user.ID.String(), 15*time.Minute); err != nil {
		return err
	}

	emailRequest := events.EmailSendRequest{
		Recipient: user.Email,
		Subject:   "Password Reset",
		Body:      fmt.Sprintf("Click here to reset your password: http://localhost:8080/reset-password?token=%s", token),
	}
	emailDataBytes, err := json.Marshal(emailRequest)
	if err != nil {
		return err
	}

	return uc.EventBus.Publish(ctx, &domain.Event{Subject: "email.send", Data: emailDataBytes})
}
