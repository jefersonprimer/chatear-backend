package notification

import (
	"context"
	"log"
	"os"

	"github.com/jefersonprimer/chatear-backend/infrastructure"
	"github.com/jefersonprimer/chatear-backend/internal/notification/application"
	"github.com/jefersonprimer/chatear-backend/internal/notification/infrastructure"
	"github.com/jefersonprimer/chatear-backend/internal/notification/worker"
	"github.com/nats-io/nats.go"
)

// Example of how to initialize and use the notification domain
func ExampleUsage() {
	// Initialize database connection
	db, err := infrastructure.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize NATS connection
	natsConn, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}

	// Initialize repositories and services
	emailRepo := notificationInfra.NewPostgresRepository(db)
	emailSender, err := notificationInfra.NewSMTPSender()
	if err != nil {
		log.Fatal("Failed to create SMTP sender:", err)
	}

	// Initialize application services
	emailService := application.NewEmailService(natsConn)
	emailUseCase := application.NewEmailSender(emailRepo, emailSender)

	// Initialize and start the NATS consumer
	consumer, err := worker.NewNatsEmailConsumer(os.Getenv("NATS_URL"), emailUseCase)
	if err != nil {
		log.Fatal("Failed to create NATS consumer:", err)
	}

	// Start the consumer in a goroutine
	ctx := context.Background()
	go consumer.Start(ctx)

	// Example: Send a welcome email
	err = emailService.SendWelcomeEmail(ctx, "user@example.com", "John Doe")
	if err != nil {
		log.Printf("Failed to send welcome email: %v", err)
	}

	// Example: Send a magic link email
	err = emailService.SendMagicLinkEmail(ctx, "user@example.com", "https://example.com/auth/magic?token=abc123")
	if err != nil {
		log.Printf("Failed to send magic link email: %v", err)
	}

	// Example: Send a custom email
	err = emailService.SendCustomEmail(ctx, "user@example.com", "Custom Subject", "Custom body content", "email")
	if err != nil {
		log.Printf("Failed to send custom email: %v", err)
	}
}