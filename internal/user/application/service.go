package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"github.com/jefersonprimer/chatear-backend/internal/user/infrastructure"
	
)

// UserApplicationService encapsulates user-related application logic.
type UserApplicationService struct {
	userRepo           domain.UserRepository
	refreshTokenRepo   domain.RefreshTokenRepository
	blacklistRepo      domain.BlacklistRepository
	eventBus           domain.EventBus
	tokenRepo          infrastructure.TokenRepository
	emailRepo          domain.EmailRepository
	tokenService       domain.TokenService
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	appURL             string
	maxEmailsPerDay    int
	userDeletionRepo   domain.UserDeletionRepository
	deletionCapacityRepo domain.DeletionCapacityRepository
}

// NewUserApplicationService creates a new UserApplicationService.
func NewUserApplicationService(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	blacklistRepo domain.BlacklistRepository,
	eventBus domain.EventBus,
	tokenRepo infrastructure.TokenRepository,
	emailRepo domain.EmailRepository,
	tokenService domain.TokenService,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
	appURL string,
	maxEmailsPerDay int,
	userDeletionRepo domain.UserDeletionRepository,
	deletionCapacityRepo domain.DeletionCapacityRepository,
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
		appURL:             appURL,
		maxEmailsPerDay:    maxEmailsPerDay,
		userDeletionRepo:   userDeletionRepo,
		deletionCapacityRepo: deletionCapacityRepo,
	}
}

func (s *UserApplicationService) Register(ctx context.Context, name, email, password string) (*AuthTokens, *domain.User, error) {
	registerUserUseCase := NewRegisterUser(s.userRepo, s.emailRepo, s.tokenRepo, s.eventBus, s.appURL, s.maxEmailsPerDay)
	user, err := registerUserUseCase.Execute(ctx, name, email, password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register user: %w", err)
	}

	accessToken, err := s.tokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenString, err := s.tokenService.CreateRefreshToken(ctx, user)
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

	if err := s.refreshTokenRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &AuthTokens{AccessToken: accessToken, RefreshToken: refreshTokenString}, user, nil
}

// Login logs in a user.
func (s *UserApplicationService) Login(ctx context.Context, email, password, ipAddress, userAgent string) (*LoginResponse, error) {
	loginUseCase := NewLoginUser(s.userRepo, s.refreshTokenRepo, s.tokenService)
	return loginUseCase.Execute(ctx, email, password, ipAddress, userAgent)
}

func (s *UserApplicationService) Logout(ctx context.Context, accessToken string, refreshToken string) error {
	// Blacklist the access token
	if accessToken != "" {
		userID, err := s.tokenService.VerifyToken(ctx, accessToken)
		if err != nil {
			// Log the error but don't fail the logout process if access token is invalid
			fmt.Printf("Warning: Failed to validate access token during logout: %v\n", err)
		} else {
			if err := s.refreshTokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
				return fmt.Errorf("failed to revoke all refresh tokens for user: %w", err)
			}
		}
	}

	return nil
}

// RecoverPassword initiates password recovery.
func (s *UserApplicationService) RecoverPassword(ctx context.Context, email string) error {
	recoverPasswordUseCase := NewPasswordRecovery(s.userRepo, s.tokenRepo, s.eventBus, s.appURL)
	return recoverPasswordUseCase.Execute(ctx, email)
}

func (s *UserApplicationService) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	deleteUserUseCase := NewDeleteUser(s.userRepo, s.userDeletionRepo, s.deletionCapacityRepo, s.eventBus)
	return deleteUserUseCase.Execute(ctx, userID)
}

// RecoverAccount recovers a user account with a token and new password.
func (s *UserApplicationService) RecoverAccount(ctx context.Context, token, newPassword string) (*AuthTokens, *domain.User, error) {
	recoverAccountUseCase := NewVerifyTokenAndResetPassword(s.userRepo, s.tokenRepo)
	user, err := recoverAccountUseCase.Execute(ctx, token, newPassword)
	if err != nil {
		return nil, nil, err
	}

	// Revoke all refresh tokens for the user
	if err := s.refreshTokenRepo.RevokeAllUserTokens(ctx, user.ID); err != nil {
		return nil, nil, fmt.Errorf("failed to revoke all refresh tokens for user: %w", err)
	}

	accessToken, err := s.tokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenString, err := s.tokenService.CreateRefreshToken(ctx, user)
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

	if err := s.refreshTokenRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &AuthTokens{AccessToken: accessToken, RefreshToken: refreshTokenString}, user, nil
}

// GetUserByID retrieves a user by their ID.
func (s *UserApplicationService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

// GetUserByEmail retrieves a user by their email.
func (s *UserApplicationService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *UserApplicationService) RefreshToken(ctx context.Context, refreshTokenString string) (*AuthTokens, *domain.User, error) {
	// Validate the refresh token
	refreshToken, err := s.refreshTokenRepo.GetRefreshTokenByToken(ctx, refreshTokenString)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if refreshToken.Revoked {
		return nil, nil, fmt.Errorf("refresh token has been revoked")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, nil, fmt.Errorf("refresh token has expired")
	}

	// Get the user associated with the refresh token
	user, err := s.userRepo.GetUserByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found for refresh token: %w", err)
	}

	// Revoke the old refresh token
	if err := s.refreshTokenRepo.RevokeRefreshToken(ctx, refreshToken.ID); err != nil {
		return nil, nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	// Generate a new access token
	newAccessToken, err := s.tokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// Generate a new refresh token
	newRefreshTokenString, err := s.tokenService.CreateRefreshToken(ctx, user)
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

	if err := s.refreshTokenRepo.CreateRefreshToken(ctx, newRefreshToken); err != nil {
		return nil, nil, fmt.Errorf("failed to save new refresh token: %w", err)
	}

	return &AuthTokens{AccessToken: newAccessToken, RefreshToken: newRefreshTokenString}, user, nil
}

// VerifyEmail verifies a user's email using a token.
func (s *UserApplicationService) VerifyEmail(ctx context.Context, token string) error {
	verifyTokenUseCase := NewVerifyToken(s.userRepo, s.tokenRepo) // Assuming s.tokenRepo exists
	return verifyTokenUseCase.Execute(ctx, token, "verification")
}

// ResendVerificationEmail resends the verification email.
func (s *UserApplicationService) ResendVerificationEmail(ctx context.Context, email string) error {
	resendVerificationEmailUseCase := NewResendVerificationEmail(s.userRepo, s.emailRepo, s.tokenRepo, s.eventBus, s.appURL, s.maxEmailsPerDay)
	return resendVerificationEmailUseCase.Execute(ctx, email)
}

// AuthTokens struct to return access and refresh tokens
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}
