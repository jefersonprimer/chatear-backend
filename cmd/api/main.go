package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/primer/chatear-backend/infrastructure"
	userInfra "github.com/primer/chatear-backend/internal/user/infrastructure"
	userDomain "github.com/primer/chatear-backend/internal/user/domain"
	userApp "github.com/primer/chatear-backend/internal/user/application"
	userHTTP "github.com/primer/chatear-backend/presentation/http"
	"github.com/primer/chatear-backend/shared/auth"
)

func main() {
	r, err := SetupServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🚀 Chatear Backend running on :8080")
	r.Run(":8080")
}

func SetupServer() (*gin.Engine, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable not set")
	}

	infra, err := infrastructure.NewInfrastructure(databaseURL, redisURL)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := userInfra.NewPostgresUserRepository(infra.Postgres)
	blacklistRepo := userInfra.NewRedisBlacklistRepository(infra.Redis)

	// Initialize services
	userService := userApp.NewService(userRepo, blacklistRepo)

	// Initialize HTTP handlers
	userHandler := userHTTP.NewUserHandler(userService)

	r := gin.Default()

	// Public routes
	publicRoutes := r.Group("/api/v1")
	{
		publicRoutes.POST("/register", userHandler.Register)
		publicRoutes.POST("/login", userHandler.Login)
	}

	// Authenticated routes
	authRoutes := r.Group("/api/v1")
	authRoutes.Use(auth.AuthMiddleware(blacklistRepo))
	{
		authRoutes.GET("/me", userHandler.GetMe)
		authRoutes.POST("/logout", userHandler.Logout)
	}

	return r, nil
}