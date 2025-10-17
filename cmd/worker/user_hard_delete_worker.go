package main

import (
	"context"
	"log"
	"time"

	"github.com/jefersonprimer/chatear-backend/config"
	"github.com/jefersonprimer/chatear-backend/infrastructure"
	userInfra "github.com/jefersonprimer/chatear-backend/internal/user/infrastructure"
)

func main() {
	cfg := config.LoadConfig()

	infra, err := infrastructure.NewInfrastructure(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.RedisURL, cfg.NatsURL)
	if err != nil {
		log.Fatal(err)
	}

	// userRepo := userInfra.NewPostgresUserRepository(infra.Postgres.Pool)
	var userRepo userInfra.UserRepository

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running hard delete worker...")
		if err := hardDeleteUsers(userRepo, cfg.HardDeleteRetentionPeriod); err != nil {
			log.Printf("Error hard deleting users: %v", err)
		}
	}
}

func hardDeleteUsers(userRepo userInfra.UserRepository, retentionPeriod time.Duration) error {
	ctx := context.Background()
	// Get users to delete
	users, err := userRepo.GetSoftDeletedUsers(ctx, retentionPeriod)
	if err != nil {
		return err
	}

	for _, user := range users {
		log.Printf("Hard deleting user %s", user.ID)
		if err := userRepo.HardDeleteUser(ctx, user.ID); err != nil {
			log.Printf("Error hard deleting user %s: %v", user.ID, err)
		}
	}

	return nil
}
