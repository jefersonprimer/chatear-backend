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
	"golang.org/x/crypto/bcrypt"
)

const maxEmailsPerDay = 2

// RegisterUser is a use case for registering a new user.
type RegisterUser struct {
	UserRepository  domain.UserRepository
	EmailRepository domain.EmailRepository
	TokenRepository infrastructure.TokenRepository
	EventBus        domain.EventBus
}

// NewRegisterUser creates a new RegisterUser use case.
func NewRegisterUser(userRepository domain.UserRepository, emailRepository domain.EmailRepository, tokenRepository infrastructure.TokenRepository, eventBus domain.EventBus) *RegisterUser {
	return &RegisterUser{
		UserRepository:  userRepository,
		EmailRepository: emailRepository,
		TokenRepository: tokenRepository,
		EventBus:        eventBus,
	}
}

// Execute registers a new user and sends a verification email.
func (uc *RegisterUser) Execute(ctx context.Context, name, email, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := uc.UserRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	emails, err := uc.EmailRepository.GetEmailsByUserIDAndType(ctx, user.ID, "verification")
	if err != nil {
		return nil, err	}

	if len(emails) >= maxEmailsPerDay {
		return nil, errors.New("email limit exceeded")
	}

	token, err := util.GenerateRandomToken()
	if err != nil {
		return nil, err
	}

	if err := uc.TokenRepository.Set(ctx, fmt.Sprintf("verification:%s", token), user.ID.String(), 15*time.Minute); err != nil {
		return nil, err
	}

	emailRequest := events.EmailSendRequest{
		Recipient: user.Email,
		Subject:   "Email Verification",
		Body:      fmt.Sprintf("Click here to verify your email: http://localhost:8080/verify-email?token=%s", token),
	}
	emailDataBytes, err := json.Marshal(emailRequest)
	if err != nil {
		return nil, err
	}

	if err := uc.EventBus.Publish(ctx, &domain.Event{Subject: "email.send", Data: emailDataBytes}); err != nil {
		return nil, err
	}

	if err := uc.EmailRepository.CreateEmail(ctx, &domain.Email{ID: uuid.New(), UserID: user.ID, Type: "verification"}); err != nil {
		return nil, err
	}

	return user, nil
}
