package postgres

import (
	"context"
	"log"
	"os"

	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// ExampleUsage demonstrates how to use the database adapter and repositories
func ExampleUsage() {
	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Create database adapter
	adapter, err := NewAdapter(databaseURL)
	if err != nil {
		log.Fatalf("Failed to create database adapter: %v", err)
	}
	defer adapter.Close()

	// Run migrations
	err = adapter.RunMigrations("./migrations/postgres")
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories
	userRepo := NewUserRepository(adapter)
	actionLogRepo := NewActionLogRepository(adapter)
	_ = NewMagicLinkRepository(adapter) // Example repository creation

	ctx := context.Background()

	// Example: Create a user
	user := entities.NewUser("John Doe", "john@example.com", "hashed_password")
	err = userRepo.Create(ctx, user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return
	}

	// Example: Create an action log
	actionLog := entities.NewActionLog(&user.ID, "user_created", map[string]any{
		"source": "api",
		"ip":     "192.168.1.1",
	})
	err = actionLogRepo.Create(ctx, actionLog)
	if err != nil {
		log.Printf("Failed to create action log: %v", err)
	}

	// Example: Retrieve user by email
	retrievedUser, err := userRepo.GetByEmail(ctx, "john@example.com")
	if err != nil {
		log.Printf("Failed to retrieve user: %v", err)
		return
	}

	log.Printf("Retrieved user: %+v", retrievedUser)
}

// ExampleRepositoryFactory creates all repositories with the given adapter
func ExampleRepositoryFactory(adapter *Adapter) map[string]interface{} {
	return map[string]interface{}{
		"user":        NewUserRepository(adapter),
		"actionLog":   NewActionLogRepository(adapter),
		"magicLink":   NewMagicLinkRepository(adapter),
		// Add other repositories as needed
	}
}
