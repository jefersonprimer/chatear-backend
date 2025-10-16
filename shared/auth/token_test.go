package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jefersonprimer/chatear-backend/shared/constants"
)

// Test constants for JWT
const (
	testJwtSecret           = "test_secret_key_for_jwt_testing_only"
	testAccessTokenDuration = time.Minute * 15 // 15 minutes for testing
)

func TestGenerateAccessToken(t *testing.T) {
	service := NewTokenService()
	userID := "test-user-id"

	// Temporarily override JwtSecret for testing
	oldJwtSecret := constants.JwtSecret
	constants.JwtSecret = []byte(testJwtSecret)
	defer func() { constants.JwtSecret = oldJwtSecret }()

	tokenString, err := service.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	if tokenString == "" {
		t.Error("Generated token string is empty")
	}

	// Validate the generated token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJwtSecret), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse generated token: %v", err)
	}

	if !token.Valid {
		t.Error("Generated token is invalid")
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	expectedExp := time.Now().Add(testAccessTokenDuration).Unix()
	if claims.ExpiresAt.Unix() < expectedExp-5 || claims.ExpiresAt.Unix() > expectedExp+5 {
		t.Errorf("Expected expiration around %d, got %d", expectedExp, claims.ExpiresAt.Unix())
	}
}

func TestValidateAccessToken(t *testing.T) {
	service := NewTokenService()
	userID := "test-user-id-2"

	// Temporarily override JwtSecret for testing
	oldJwtSecret := constants.JwtSecret
	constants.JwtSecret = []byte(testJwtSecret)
	defer func() { constants.JwtSecret = oldJwtSecret }()

	// Generate a valid token
	expirationTime := time.Now().Add(testAccessTokenDuration)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validTokenString, err := token.SignedString([]byte(testJwtSecret))
	if err != nil {
		t.Fatalf("Failed to sign token for validation test: %v", err)
	}

	// Test case 1: Valid token
	validatedClaims, err := service.ValidateAccessToken(validTokenString)
	if err != nil {
		t.Fatalf("ValidateAccessToken failed for valid token: %v", err)
	}
	if validatedClaims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, validatedClaims.UserID)
	}

	// Test case 2: Invalid token (wrong secret)
	invalidToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("wrong_secret"))
	_, err = service.ValidateAccessToken(invalidToken)
	if err == nil {
		t.Error("ValidateAccessToken unexpectedly succeeded for invalid token")
	}

	// Test case 3: Expired token
	expiredClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
		},
	}
	expiredToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte(testJwtSecret))
	_, err = service.ValidateAccessToken(expiredToken)
	if err == nil {
		t.Fatal("ValidateAccessToken unexpectedly succeeded for expired token")
	}
	// Check if the error indicates token expiration
	if err != nil && err.Error() != "failed to parse access token: Token is expired" {
		t.Errorf("Expected expired token error, got %v", err)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	service := NewTokenService()
	token, err := service.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if token == "" {
		t.Error("Generated refresh token is empty")
	}
	// Basic check for format, acknowledging it's a placeholder
	if len(token) < 10 {
		t.Errorf("Generated refresh token is too short: %s", token)
	}
}

func TestValidateRefreshToken(t *testing.T) {
	service := NewTokenService()

	// Test case 1: Valid placeholder token
	valid, err := service.ValidateRefreshToken("any-non-empty-string")
	if err != nil {
		t.Fatalf("ValidateRefreshToken failed for valid placeholder: %v", err)
	}
	if !valid {
		t.Error("Expected valid for placeholder token, got false")
	}

	// Test case 2: Empty token
	valid, err = service.ValidateRefreshToken("")
	if err == nil {
		t.Error("ValidateRefreshToken unexpectedly succeeded for empty token")
	}
	if valid {
		t.Error("Expected invalid for empty token, got true")
	}
	if err != nil && err.Error() != "refresh token is empty" {
		t.Errorf("Expected 'refresh token is empty' error, got %v", err)
	}
}