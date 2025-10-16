package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
}

// NewTokenService creates a new TokenService
func NewTokenService(refreshTokenRepo domain.RefreshTokenRepository) *TokenService {
	return &TokenService{refreshTokenRepo: refreshTokenRepo}
}

// GenerateAccessToken creates a new JWT access token
func (s *TokenService) GenerateAccessToken(userID string) (string, error) {
	expirationTime := time.Now().Add(constants.AccessTokenExpiration)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(constants.JwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return tokenString, nil
}

// ValidateAccessToken validates the JWT access token and returns the claims
func (s *TokenService) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return constants.JwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	return claims, nil
}

// GenerateRefreshToken creates a new refresh token (UUID for now, will be stored in DB)
func (s *TokenService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes for refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateRefreshToken validates a refresh token (this will involve DB lookup later)
func (s *TokenService) ValidateRefreshToken(ctx context.Context, tokenString string) (*domain.RefreshToken, error) {
	refreshToken, err := s.refreshTokenRepo.FindByToken(tokenString)
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