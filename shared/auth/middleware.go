package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	userDomain "github.com/primer/chatear-backend/internal/user/domain"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
)

// AuthMiddleware creates a Gin middleware for JWT authentication.
func AuthMiddleware(blacklistRepo userDomain.BlacklistRepository) gin.HandlerFunc {
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

		claims, err := ValidateAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		isBlacklisted, err := blacklistRepo.Check(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check token blacklist"})
			return
		}
		if isBlacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		c.Set(string(ContextKeyUserID), userID)
		c.Next()
	}
}

// GetUserIDFromContext extracts the UserID from the context.
func GetUserIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
	userID, ok := ctx.Get(string(ContextKeyUserID))
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return userID.(uuid.UUID), nil
}