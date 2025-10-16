package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"github.com/jefersonprimer/chatear-backend/internal/user/infrastructure"
)

// VerifyToken is a use case for verifying a token.
type VerifyToken struct {
	UserRepository  domain.UserRepository
	TokenRepository infrastructure.TokenRepository
}

// NewVerifyToken creates a new VerifyToken use case.
func NewVerifyToken(userRepository domain.UserRepository, tokenRepository infrastructure.TokenRepository) *VerifyToken {
	return &VerifyToken{
		UserRepository:  userRepository,
		TokenRepository: tokenRepository,
	}
}

// Execute verifies a token and performs the corresponding action.
func (uc *VerifyToken) Execute(ctx context.Context, token, tokenType string) error {
	key := fmt.Sprintf("%s:%s", tokenType, token)
	userIDString, err := uc.TokenRepository.Get(ctx, key)
	if err != nil {
		return err
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return err
	}

	user, err := uc.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	switch tokenType {
	case "verification":
		user.IsEmailVerified = true
		if err := uc.UserRepository.UpdateUser(ctx, user); err != nil {
			return err
		}
	case "password-reset":
		// In a real application, you would allow the user to set a new password.
		// For now, we'll just log that the password reset was successful.
		fmt.Printf("Password reset for user %s was successful\n", user.Email)
	}

	return uc.TokenRepository.Del(ctx, key)
}
