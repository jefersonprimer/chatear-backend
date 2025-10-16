package application

import (
	"context"
	"time"

	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
)

// LogoutUser is a use case for logging out a user.
type LogoutUser struct {
	RefreshTokenRepository domain.RefreshTokenRepository
}

// NewLogoutUser creates a new LogoutUser use case.
func NewLogoutUser(refreshTokenRepository domain.RefreshTokenRepository) *LogoutUser {
	return &LogoutUser{RefreshTokenRepository: refreshTokenRepository}
}

// Execute logs out a user by invalidating their refresh token.
func (uc *LogoutUser) Execute(ctx context.Context, token string) error {
	refreshToken, err := uc.RefreshTokenRepository.GetRefreshTokenByToken(ctx, token)
	if err != nil {
		return err
	}

	refreshToken.Revoked = true
	now := time.Now()
	refreshToken.RevokedAt = &now
	return uc.RefreshTokenRepository.UpdateRefreshToken(ctx, refreshToken)
}
