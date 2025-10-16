package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/primer/chatear-backend/shared/constants"
)

// Claims defines the structure of our JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a new JWT access token
func GenerateAccessToken(userID string) (string, error) {
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
func ValidateAccessToken(tokenString string) (*Claims, error) {
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
func GenerateRefreshToken() (string, error) {
	// For now, a simple UUID-like string. This will be stored in the database.
	// In a real application, you might use a cryptographically secure random string.
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix()), nil
}

// ValidateRefreshToken validates a refresh token (this will involve DB lookup later)
func ValidateRefreshToken(token string) (bool, error) {
	// Placeholder for now. This will involve checking against the database.
	if token == "" {
		return false, fmt.Errorf("refresh token is empty")
	}
	return true, nil
}