package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/application"
	"github.com/jefersonprimer/chatear-backend/shared/auth"
)

// UserHandlers handles HTTP requests for user operations
type UserHandlers struct {
	userService *application.UserApplicationService
}

// NewUserHandlers creates a new user handlers instance
func NewUserHandlers(userService *application.UserApplicationService) *UserHandlers {
	return &UserHandlers{
		userService: userService,
	}
}

// Register handles POST /register
func (h *UserHandlers) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authTokens, user, err := h.userService.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "User registered successfully",
		"access_token":  authTokens.AccessToken,
		"refresh_token": authTokens.RefreshToken,
		"user":          user,
	})
}

// Login handles POST /login
func (h *UserHandlers) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authTokens, user, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Logged in successfully",
		"access_token":  authTokens.AccessToken,
		"refresh_token": authTokens.RefreshToken,
		"user":          user,
	})
}

// Logout handles POST /logout
func (h *UserHandlers) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	accessToken := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		accessToken = authHeader[7:]
	}

	if err := h.userService.Logout(c.Request.Context(), accessToken, req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken handles POST /refresh-token
func (h *UserHandlers) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authTokens, user, err := h.userService.RefreshToken(c.Request.Context(), req.RefreshToken)
		if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Token refreshed successfully",
		"access_token":  authTokens.AccessToken,
		"refresh_token": authTokens.RefreshToken,
		"user":          user,
	})
}

// VerifyEmail handles GET /verify-email
func (h *UserHandlers) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	if err := h.userService.VerifyEmail(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// GetMe handles GET /me
func (h *UserHandlers) GetMe(c *gin.Context) {
	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}