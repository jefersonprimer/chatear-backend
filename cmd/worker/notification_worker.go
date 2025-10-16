package main

import (
	"context"
	"log"
	"os"

	"github.com/jefersonprimer/chatear-backend/infrastructure"
	notification_app "github.com/jefersonprimer/chatear-backend/internal/notification/application"
	notification_infra "github.com/jefersonprimer/chatear-backend/internal/notification/infrastructure"
	"github.com/jefersonprimer/chatear-backend/internal/notification/worker"
	"github.com/joho/godotenv"
)

type EmailSendEvent struct {
	Recipient    string `json:"recipient"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	TemplateName string `json:"template_name,omitempty"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	db, err := infrastructure.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	notificationRepository := notification_infra.NewPostgresRepository(db)
	smtpSender, err := notification_infra.NewSMTPSender()
	if err != nil {
		log.Fatalf("Error creating SMTP sender: %v", err)
	}
	emailSender := notification_app.NewEmailSender(notificationRepository, smtpSender)

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL environment variable not set")
	}

	natsConsumer, err := worker.NewNatsEmailConsumer(natsURL, emailSender)
	if err != nil {
		log.Fatalf("Error creating NATS consumer: %v", err)
	}

	ctx := context.Background()
	log.Println("Starting notification worker...")
	natsConsumer.Start(ctx)
}