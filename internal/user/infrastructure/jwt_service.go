package infrastructure

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
)

// JWTService is a JWT implementation of the domain.TokenService.
type JWTService struct {
	SecretKey []byte
}

// NewJWTService creates a new JWTService.
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{SecretKey: []byte(secretKey)}
}

// CreateAccessToken creates a new access token for the given user.
func (s *JWTService) CreateAccessToken(ctx context.Context, user *domain.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.SecretKey)
}

// CreateRefreshToken creates a new refresh token for the given user.
func (s *JWTService) CreateRefreshToken(ctx context.Context, user *domain.User) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// VerifyToken verifies the given token and returns the user ID.
func (s *JWTService) VerifyToken(ctx context.Context, tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.SecretKey, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return uuid.Parse(claims.Subject)
	} else {
		return uuid.Nil, err
	}
}
