package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	userDomain "github.com/jefersonprimer/chatear-backend/internal/user/domain"
)

type contextKey string

const (
	ContextKeyUserID       contextKey = "userID"
	ContextKeyRefreshToken contextKey = "refreshToken"
	ContextKeyAccessToken  contextKey = "accessToken"
)

// AuthMiddleware creates a Gin middleware for JWT authentication.
func AuthMiddleware(tokenService *TokenService, blacklistRepo userDomain.BlacklistRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := authHeader
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		userID, err := tokenService.VerifyToken(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Store userID in Gin context
		c.Set(string(ContextKeyUserID), userID)

		// Extract refresh token from header
		refreshToken := c.GetHeader("X-Refresh-Token")

		// Store userID, accessToken, and refreshToken in request context for GraphQL resolvers
		ctx := context.WithValue(c.Request.Context(), ContextKeyUserID, userID)
		ctx = context.WithValue(ctx, ContextKeyAccessToken, tokenString)
		ctx = context.WithValue(ctx, ContextKeyRefreshToken, refreshToken)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetUserIDFromContext extracts the UserID from the context.
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(ContextKeyUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// GetAccessTokenFromContext extracts the AccessToken from the context.
func GetAccessTokenFromContext(ctx context.Context) (string, error) {
	accessToken, ok := ctx.Value(ContextKeyAccessToken).(string)
	if !ok {
		return "", fmt.Errorf("access token not found in context")
	}
	return accessToken, nil
}

// GetRefreshTokenFromContext extracts the RefreshToken from the context.
func GetRefreshTokenFromContext(ctx context.Context) (string, error) {
	refreshToken, ok := ctx.Value(ContextKeyRefreshToken).(string)
	if !ok {
		return "", fmt.Errorf("refresh token not found in context")
	}
	return refreshToken, nil
}