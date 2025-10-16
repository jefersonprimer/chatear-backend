package application

import (
	"context"
	"fmt"
	"time"

	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"github.com/jefersonprimer/chatear-backend/shared/auth"
)

// UserApplicationService encapsulates user-related application logic.
type UserApplicationService struct {
	userRepo           domain.UserRepository
	refreshTokenRepo   domain.RefreshTokenRepository
	blacklistRepo      domain.BlacklistRepository
	eventBus           domain.EventBus
	tokenRepo          infrastructure.TokenRepository
	emailRepo          domain.EmailRepository
	tokenService       *auth.TokenService
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewUserApplicationService creates a new UserApplicationService.
func NewUserApplicationService(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	blacklistRepo domain.BlacklistRepository,
	eventBus domain.EventBus,
	tokenRepo infrastructure.TokenRepository,
	emailRepo domain.EmailRepository,
	tokenService *auth.TokenService,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) *UserApplicationService {
	return &UserApplicationService{
		userRepo:           userRepo,
		refreshTokenRepo:   refreshTokenRepo,
		blacklistRepo:      blacklistRepo,
		eventBus:           eventBus,
		tokenRepo:          tokenRepo,
		emailRepo:          emailRepo,
		tokenService:       tokenService,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// Register registers a new user.
func (s *UserApplicationService) Register(ctx context.Context, name, email, password string) (*AuthTokens, *domain.User, error) {
	registerUserUseCase := NewRegisterUser(s.userRepo, s.emailRepo, s.tokenRepo, s.eventBus)
	user, err := registerUserUseCase.Execute(ctx, name, email, password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register user: %w", err)
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenString, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token string: %w", err)
	}

	refreshToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(s.refreshTokenDuration),
		CreatedAt: time.Now(),
	}

	if err := s.refreshTokenRepo.Save(refreshToken); err != nil {
		return nil, nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &AuthTokens{AccessToken: accessToken, RefreshToken: refreshTokenString}, user, nil
}

// Login logs in a user.
func (s *UserApplicationService) Login(ctx context.Context, email, password string) (*AuthTokens, *domain.User, error) {
	loginUseCase := NewLogin(s.userRepo)
	user, err := loginUseCase.Execute(ctx, email, password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to login user: %w", err)
	}

	accessToken, err := auth.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenString, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token string: %w", err)
	}

	refreshToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(s.refreshTokenDuration),
		CreatedAt: time.Now(),
	}

	if err := s.refreshTokenRepo.Save(refreshToken); err != nil {
		return nil, nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &AuthTokens{AccessToken: accessToken, RefreshToken: refreshTokenString}, user, nil
}

// Logout logs out a user.
func (s *UserApplicationService) Logout(ctx context.Context, accessToken string, refreshToken string) error {
	// Blacklist the access token
	if accessToken != "" {
		claims, err := auth.ValidateAccessToken(accessToken)
		if err != nil {
			// Log the error but don't fail the logout process if access token is invalid
			fmt.Printf("Warning: Failed to validate access token during logout: %v\n", err)
		} else {
			expiration := claims.ExpiresAt.Sub(time.Now())
			if expiration > 0 {
				err := s.blacklistRepo.Add(ctx, accessToken, expiration)
				if err != nil {
					return fmt.Errorf("failed to blacklist access token: %w", err)
				}
			}
		}
	}

	// Find the refresh token in the database
	dbRefreshToken, err := s.refreshTokenRepo.FindByToken(refreshToken)
	if err != nil {
		return fmt.Errorf("failed to find refresh token for logout: %w", err)
	}

	// Revoke the refresh token
	if err := s.refreshTokenRepo.Revoke(dbRefreshToken.ID); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

// RecoverPassword initiates password recovery.
func (s *UserApplicationService) RecoverPassword(ctx context.Context, email string) error {
	recoverPasswordUseCase := NewRecoverPassword(s.userRepo, s.eventBus)
	return recoverPasswordUseCase.Execute(ctx, email)
}

// DeleteAccount deletes a user account.
func (s *UserApplicationService) DeleteAccount(ctx context.Context, userID uuid.UUID, password string) error {
	deleteUserUseCase := NewDeleteUser(s.userRepo, s.eventBus)
	return deleteUserUseCase.Execute(ctx, userID, password)
}

// RecoverAccount recovers a user account with a token and new password.
func (s *UserApplicationService) RecoverAccount(ctx context.Context, token, newPassword string) (*AuthTokens, *domain.User, error) {
	recoverAccountUseCase := NewVerifyTokenAndResetPassword(s.userRepo, s.eventBus)
	user, err := recoverAccountUseCase.Execute(ctx, token, newPassword)
	if err != nil {
		return nil, nil, err
	}

	accessToken, err := auth.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenString, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token string: %w", err)
	}

	refreshToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(s.refreshTokenDuration),
		CreatedAt: time.Now(),
	}

	if err := s.refreshTokenRepo.Save(refreshToken); err != nil {
		return nil, nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &AuthTokens{AccessToken: accessToken, RefreshToken: refreshTokenString}, user, nil
}

// GetUserByID retrieves a user by their ID.
func (s *UserApplicationService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *UserApplicationService) RefreshToken(ctx context.Context, refreshTokenString string) (*AuthTokens, *domain.User, error) {
	// Validate the refresh token
	refreshToken, err := s.tokenService.ValidateRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get the user associated with the refresh token
	user, err := s.userRepo.GetUserByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found for refresh token: %w", err)
	}

	// Revoke the old refresh token
	if err := s.refreshTokenRepo.Revoke(refreshToken.ID); err != nil {
		return nil, nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	// Generate a new access token
	newAccessToken, err := s.tokenService.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// Generate a new refresh token
	newRefreshTokenString, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate new refresh token string: %w", err)
	}

	newRefreshToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     newRefreshTokenString,
		ExpiresAt: time.Now().Add(s.refreshTokenDuration),
		CreatedAt: time.Now(),
	}

	if err := s.refreshTokenRepo.Save(newRefreshToken); err != nil {
		return nil, nil, fmt.Errorf("failed to save new refresh token: %w", err)
	}

	return &AuthTokens{AccessToken: newAccessToken, RefreshToken: newRefreshTokenString}, user, nil
}

// VerifyEmail verifies a user's email using a token.
func (s *UserApplicationService) VerifyEmail(ctx context.Context, token string) error {
	verifyTokenUseCase := NewVerifyToken(s.userRepo, s.tokenRepo) // Assuming s.tokenRepo exists
	return verifyTokenUseCase.Execute(ctx, token, "verification")
}

// AuthTokens struct to return access and refresh tokens
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}
