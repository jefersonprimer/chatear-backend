package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"github.com/jefersonprimer/chatear-backend/internal/user/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

// VerifyTokenAndResetPassword is a use case for verifying a token and resetting a password.
type VerifyTokenAndResetPassword struct {
	UserRepository  domain.UserRepository
	TokenRepository infrastructure.TokenRepository
}

// NewVerifyTokenAndResetPassword creates a new VerifyTokenAndResetPassword use case.
func NewVerifyTokenAndResetPassword(userRepository domain.UserRepository, tokenRepository infrastructure.TokenRepository) *VerifyTokenAndResetPassword {
	return &VerifyTokenAndResetPassword{
		UserRepository:  userRepository,
		TokenRepository: tokenRepository,
	}
}

// Execute verifies a token, resets the password, and returns the user.
func (uc *VerifyTokenAndResetPassword) Execute(ctx context.Context, token, newPassword string) (*domain.User, error) {
	userIDString, err := uc.TokenRepository.Get(ctx, fmt.Sprintf("password-reset:%s", token))
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	user, err := uc.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = string(hashedPassword)
	if err := uc.UserRepository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	if err := uc.TokenRepository.Del(ctx, fmt.Sprintf("password-reset:%s", token)); err != nil {
		// Log the error but don't fail the process
		fmt.Printf("Warning: Failed to delete password reset token: %v\n", err)
	}

	return user, nil
}
