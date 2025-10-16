package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/application"
)

// UserHandler holds the use cases for the user domain.
type UserHandler struct {
	RegisterUser     *application.RegisterUser
	LoginUser        *application.LoginUser
	LogoutUser       *application.LogoutUser
	PasswordRecovery *application.PasswordRecovery
	DeleteUser       *application.DeleteUser
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(registerUser *application.RegisterUser, loginUser *application.LoginUser, logoutUser *application.LogoutUser, passwordRecovery *application.PasswordRecovery, deleteUser *application.DeleteUser) *UserHandler {
	return &UserHandler{
		RegisterUser:     registerUser,
		LoginUser:        loginUser,
		LogoutUser:       logoutUser,
		PasswordRecovery: passwordRecovery,
		DeleteUser:       deleteUser,
	}
}

// Register is the handler for the user registration endpoint.
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.RegisterUser.Execute(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login is the handler for the user login endpoint.
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.LoginUser.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout is the handler for the user logout endpoint.
func (h *UserHandler) Logout(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.LogoutUser.Execute(c.Request.Context(), req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// PasswordRecovery is the handler for the password recovery endpoint.
func (h *UserHandler) PasswordRecovery(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.PasswordRecovery.Execute(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password recovery email sent"})
}

// Delete is the handler for the user deletion endpoint.
func (h *UserHandler) Delete(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.DeleteUser.Execute(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}