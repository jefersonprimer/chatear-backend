package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser is a use case for logging in a user.
type LoginUser struct {
	UserRepository         domain.UserRepository
	RefreshTokenRepository domain.RefreshTokenRepository
	TokenService           domain.TokenService
}

// NewLoginUser creates a new LoginUser use case.
func NewLoginUser(userRepository domain.UserRepository, refreshTokenRepository domain.RefreshTokenRepository, tokenService domain.TokenService) *LoginUser {
	return &LoginUser{
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
		TokenService:           tokenService,
	}
}

// LoginResponse is the response for the login use case.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Execute logs in a user and returns an access token and a refresh token.
func (uc *LoginUser) Execute(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := uc.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, err
	}

	accessToken, err := uc.TokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.TokenService.CreateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}

	refreshTokenEntity := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
		Revoked:   false,
	}

	if err := uc.RefreshTokenRepository.CreateRefreshToken(ctx, refreshTokenEntity); err != nil {
		return nil, err
	}

	return &LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
