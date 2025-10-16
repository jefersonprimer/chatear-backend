package main

import (
	"context"
	"log"
	"os"

	"github.com/jefersonprimer/chatear-backend/infrastructure"
	notification_app "github.com/jefersonprimer/chatear-backend/internal/notification/application"
	notification_infra "github.com/jefersonprimer/chatear-backend/internal/notification/infrastructure"
	"github.com/jefersonprimer/chatear-backend/internal/notification/worker"
	"github.com/jefersonprimer/chatear-backend/shared/events"
	"github.com/joho/godotenv"
)
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

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
		log.Fatalf("Error initializing infrastructure: %v", err)
	}
	defer infra.Close()

	notificationRepository := notification_infra.NewPostgresRepository(infra.Postgres)
	smtpSender, err := notification_infra.NewSMTPSender()
	if err != nil {
		log.Fatalf("Error creating SMTP sender: %v", err)
	}
	emailSender := notification_app.NewEmailSender(notificationRepository, smtpSender)

	natsConsumer, err := worker.NewNatsEmailConsumer(infra.NatsConn, emailSender)
	if err != nil {
		log.Fatalf("Error creating NATS consumer: %v", err)
	}

	ctx := context.Background()
	log.Println("Starting notification worker...")
	natsConsumer.Start(ctx)
}