package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/jefersonprimer/chatear-backend/graph"
	"github.com/jefersonprimer/chatear-backend/infrastructure"
	userInfra "github.com/jefersonprimer/chatear-backend/internal/user/infrastructure"
	userDomain "github.com/jefersonprimer/chatear-backend/internal/user/domain"
	userApp "github.com/jefersonprimer/chatear-backend/internal/user/application"
	userHTTP "github.com/jefersonprimer/chatear-backend/presentation/http"
	"github.com/jefersonprimer/chatear-backend/shared/auth"
	"github.com/jefersonprimer/chatear-backend/shared/constants"
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

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL environment variable not set")
	}

	infra, err := infrastructure.NewInfrastructure(databaseURL, redisURL, natsURL)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := userInfra.NewPostgresUserRepository(infra.Postgres.Pool)
	blacklistRepo := userInfra.NewRedisBlacklistRepository(infra.Redis)
	refreshTokenRepo := userInfra.NewPostgresRefreshTokenRepository(infra.Postgres)
	emailRepo := userInfra.NewPostgresEmailRepository(infra.Postgres)
	tokenRepo := userInfra.NewRedisTokenRepository(infra.Redis)
	deletionCapacityRepo := userInfra.NewPostgresDeletionCapacityRepository(infra.Postgres)
	userDeletionRepo := userInfra.NewPostgresUserDeletionRepository(infra.Postgres)

	// Initialize event bus (NATS for example)
	eventBus := userInfra.NewNatsEventBus(infra.NatsConn)

	// Initialize shared services
	tokenService := auth.NewTokenService(refreshTokenRepo)

	// Initialize user application services
	userAppService := userApp.NewUserApplicationService(
		userRepo,
		refreshTokenRepo,
		blacklistRepo,
		eventBus,
		tokenRepo,
		emailRepo,
		tokenService,
		constants.AccessTokenExpiration,
		constants.RefreshTokenExpiration,
	)

	// Initialize HTTP handlers
	userHandler := userHTTP.NewUserHandler(userAppService)

	r := gin.Default()

	// Public routes
	publicRoutes := r.Group("/api/v1")
	{
		publicRoutes.POST("/register", userHandler.Register)
		publicRoutes.POST("/login", userHandler.Login)
		publicRoutes.GET("/verify-email", userHandler.VerifyEmail)
		publicRoutes.POST("/refresh-token", userHandler.RefreshToken)
	}

	// Authenticated routes
	authRoutes := r.Group("/api/v1")
	authRoutes.Use(auth.AuthMiddleware(tokenService, blacklistRepo))
	{
		authRoutes.GET("/me", userHandler.GetMe)
		authRoutes.POST("/logout", userHandler.Logout)
	}

	// GraphQL setup
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			UserAppService: userAppService,
			TokenService:   tokenService,
		},
	}))

	r.POST("/graphql", func(c *gin.Context) {
		// Apply JWT middleware to GraphQL endpoint
		// This middleware will extract the token and add user info to context
		auth.AuthMiddleware(blacklistRepo)(c)
		// Pass the Gin context to the GraphQL handler
		srv.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/playground", playground.Handler("GraphQL playground", "/graphql"))

	return r, nil
}