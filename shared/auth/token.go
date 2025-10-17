package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"github.com/jefersonprimer/chatear-backend/shared/constants"
)

// Claims defines the structure of our JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenService struct {
	refreshTokenRepo domain.RefreshTokenRepository
	jwtSecret        []byte
}

// NewTokenService creates a new TokenService
func NewTokenService(refreshTokenRepo domain.RefreshTokenRepository, jwtSecret string) *TokenService {
	return &TokenService{
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        []byte(jwtSecret),
	}
}

func (s *TokenService) CreateAccessToken(ctx context.Context, user *domain.User) (string, error) {
	expirationTime := time.Now().Add(constants.AccessTokenExpiration)
	claims := &Claims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return tokenString, nil
}

func (s *TokenService) VerifyToken(ctx context.Context, tokenString string) (uuid.UUID, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	return userID, nil
}

func (s *TokenService) CreateRefreshToken(ctx context.Context, user *domain.User) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes for refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateRefreshToken validates a refresh token (this will involve DB lookup later)
func (s *TokenService) ValidateRefreshToken(ctx context.Context, tokenString string) (*domain.RefreshToken, error) {
	refreshToken, err := s.refreshTokenRepo.GetRefreshTokenByToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or invalid: %w", err)
	}

	if refreshToken.Revoked {
		return nil, fmt.Errorf("refresh token has been revoked")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token has expired")
	}

	return refreshToken, nil
}