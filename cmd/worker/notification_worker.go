package main

import (
	"context"
	"log"

	"github.com/jefersonprimer/chatear-backend/config"
	"github.com/jefersonprimer/chatear-backend/infrastructure"
	notification_app "github.com/jefersonprimer/chatear-backend/internal/notification/application"
	notification_infra "github.com/jefersonprimer/chatear-backend/internal/notification/infrastructure"
	"github.com/jefersonprimer/chatear-backend/internal/notification/worker"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.RedisURL == "" {
		log.Fatal("REDIS_URL environment variable not set")
	}

	if cfg.NatsURL == "" {
		log.Fatal("NATS_URL environment variable not set")
	}

	infra, err := infrastructure.NewInfrastructure(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.RedisURL, cfg.NatsURL)
	if err != nil {
		log.Fatalf("Error initializing infrastructure: %v", err)
	}
	defer infra.Close()

	// notificationRepository := notification_infra.NewPostgresRepository(infra.Postgres)
	var notificationRepository notification_app.EmailSendRepository
	smtpSender, err := notification_infra.NewSMTPSender(cfg)
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
