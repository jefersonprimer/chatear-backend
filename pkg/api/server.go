package api

import (
	

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/jefersonprimer/chatear-backend/config"
	"github.com/jefersonprimer/chatear-backend/graph"
	"github.com/jefersonprimer/chatear-backend/infrastructure"
	userApp "github.com/jefersonprimer/chatear-backend/internal/user/application"
	userInfra "github.com/jefersonprimer/chatear-backend/internal/user/infrastructure"
	userHTTP "github.com/jefersonprimer/chatear-backend/presentation/http"
	"github.com/jefersonprimer/chatear-backend/presentation/middleware"
	"github.com/jefersonprimer/chatear-backend/shared/auth"
)

func SetupServer(cfg *config.Config) (*gin.Engine, error) {

	infra, err := infrastructure.NewInfrastructure(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.RedisURL, cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	var userRepo userApp.UserRepository
	blacklistRepo := userInfra.NewRedisBlacklistRepository(infra.Redis)
	var refreshTokenRepo userApp.RefreshTokenRepository
	var emailRepo userApp.EmailRepository
	tokenRepo, err := userInfra.NewTokenCache(cfg.RedisURL)
	if err != nil {
		return nil, err
	}
	var deletionCapacityRepo userApp.DeletionCapacityRepository
	var userDeletionRepo userApp.UserDeletionRepository

	// Initialize event bus (NATS for example)
	eventBus := userInfra.NewNATSEventBus(infra.NatsConn)

	// Initialize shared services
	tokenService := auth.NewTokenService(refreshTokenRepo, cfg.JwtSecret)

	// Initialize user application services
	userAppService := userApp.NewUserApplicationService(
		userRepo,
		refreshTokenRepo,
		blacklistRepo,
		eventBus,
		tokenRepo,
		emailRepo,
		tokenService,
		cfg.AccessTokenTTL,
		cfg.RefreshTokenTTL,
		cfg.AppURL,
		cfg.MaxEmailsPerDay,
		userDeletionRepo,
		deletionCapacityRepo,
	)

	// Initialize HTTP handlers
	userHandler := userHTTP.NewUserHandlers(userAppService)

	r := gin.Default()

	// Public routes
	publicRoutes := r.Group("/api/v1")
	{
		publicRoutes.POST("/register", userHandler.Register)
		publicRoutes.POST("/login", userHandler.Login)
		publicRoutes.GET("/verify-email", userHandler.VerifyEmail)
		publicRoutes.POST("/refresh-token", userHandler.RefreshToken)
		publicRoutes.POST("/resend-verification-email", userHandler.ResendVerificationEmail)

		// Health check routes
		healthHandler := userHTTP.NewHealthHandler(infra, cfg)
		publicRoutes.GET("/healthz", healthHandler.Healthz)
		publicRoutes.GET("/readyz", healthHandler.Readyz)
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

	graphqlHandler := gin.WrapH(srv)
	r.POST("/graphql", middleware.GinContextToContextMiddleware(), graphqlHandler)

	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))

	return r, nil
}