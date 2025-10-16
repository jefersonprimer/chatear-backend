package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/jefersonprimer/chatear-backend/infrastructure"
	notification_app "github.com/jefersonprimer/chatear-backend/internal/notification/application"
	notification_infra "github.com/jefersonprimer/chatear-backend/internal/notification/infrastructure"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/joho/godotenv"
)

const (
	natsURL         = "NATS_URL" // Replace with actual environment variable or config
	subjectUserDelete = "user.delete"
	deletionWarningPeriod = 24 * time.Hour // 24 hours before actual deletion

	// Redis rate limit keys and limits
	redisKeyGlobalDeletionCount = "global:deletion:count:" // Suffix with YYYY-MM-DD
	redisKeyUserEmailCount      = "user:email:count:"      // Suffix with UserID:YYYY-MM-DD
	maxGlobalDeletionsPerDay    = 10
	maxEmailsPerUserPerDay      = 2
)

// UserDeleteEvent represents the payload of a user.delete event
type UserDeleteEvent struct {
	UserID string `json:"user_id"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	// Initialize NATS connection
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL environment variable not set")
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Printf("Connected to NATS at %s", natsURL)

	// Initialize database connection
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}
	db, err := infrastructure.NewDB(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL database")

	// Initialize Redis client
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable not set")
	}
	redisOpt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	rc := redis.NewClient(redisOpt)
	defer rc.Close()

	_, err = rc.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	// Initialize notification services for sending recovery emails
	notificationRepository := notification_infra.NewPostgresRepository(db)
	smtpSender, err := notification_infra.NewSMTPSender()
	if err != nil {
		log.Fatalf("Error creating SMTP sender: %v", err)
	}
	emailSender := notification_app.NewEmailSender(notificationRepository, smtpSender)

	// Subscribe to user.delete events
	_, err = nc.Subscribe(subjectUserDelete, func(msg *nats.Msg) {
		log.Printf("Received message on subject: %s", msg.Subject)
		var event UserDeleteEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error unmarshaling user.delete event: %v", err)
			return
		}
		log.Printf("Processing user deletion request for UserID: %s", event.UserID)
		handleUserDeletionEvent(context.Background(), db, rc, emailSender, event.UserID)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to %s: %v", subjectUserDelete, err)
	}

	log.Printf("Subscribed to NATS subject: %s", subjectUserDelete)

	// Start the periodic deletion checker
	go startDeletionChecker(context.Background(), db, rc, emailSender)

	log.Println("User deletion worker started")

	// Keep the worker running
	select {}
}

func handleUserDeletionEvent(ctx context.Context, db *infrastructure.DB, rc *redis.Client, emailSender *notification_app.EmailSender, userID string) {
	log.Printf("User deletion event received for UserID: %s", userID)

	// Insert into user_deletions table
	scheduledDate := time.Now().Add(24 * time.Hour) // Schedule deletion 24 hours from now
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO user_deletions (user_id, scheduled_date, status, created_at) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET 
			scheduled_date = EXCLUDED.scheduled_date,
			status = EXCLUDED.status,
			created_at = EXCLUDED.created_at
	`, userID, scheduledDate, "queued", time.Now())
	
	if err != nil {
		log.Printf("Error inserting user deletion record for UserID %s: %v", userID, err)
		return
	}

	log.Printf("User deletion scheduled for UserID: %s at %s", userID, scheduledDate.Format(time.RFC3339))
}

func startDeletionChecker(ctx context.Context, db *infrastructure.DB, rc *redis.Client, emailSender *notification_app.EmailSender) {
	ticker := time.NewTicker(1 * time.Hour) // Check every hour
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Println("Scanning user_deletions table for pending actions...")
			checkAndDeleteUsers(ctx, db, rc, emailSender)
		}
	}
}

func checkAndDeleteUsers(ctx context.Context, db *infrastructure.DB, rc *redis.Client, emailSender *notification_app.EmailSender) {
	// Query user_deletions table for pending deletions
	rows, err := db.Pool.Query(ctx, `
		SELECT id, user_id, scheduled_date, status 
		FROM user_deletions 
		WHERE status IN ('queued', 'scheduled') 
		AND scheduled_date <= $1
	`, time.Now().Add(deletionWarningPeriod))
	
	if err != nil {
		log.Printf("Error querying user_deletions table: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, userID string
		var scheduledDate time.Time
		var status string
		if err := rows.Scan(&id, &userID, &scheduledDate, &status); err != nil {
			log.Printf("Error scanning user_deletions row: %v", err)
			continue
		}

		now := time.Now()
		
		// Check if we should send recovery email (24 hours before deletion)
		if status == "queued" && now.Add(deletionWarningPeriod).After(scheduledDate) {
			// Send recovery email if within 24 hours of deletion
			if !checkUserEmailRateLimit(ctx, rc, userID) {
				log.Printf("Skipping recovery email for UserID %s: rate limit exceeded", userID)
				continue
			}
			
			log.Printf("Sending recovery email to UserID: %s. Deletion scheduled for: %s", userID, scheduledDate.Format(time.RFC3339))
			
			// Send recovery email
			_, err := emailSender.Send(ctx, userID, "Account Deletion Recovery", 
				"Your account is scheduled for deletion. Click here to recover your account.", "recovery")
			if err != nil {
				log.Printf("Error sending recovery email to UserID %s: %v", userID, err)
				continue
			}
			
			// Update status to scheduled
			_, err = db.Pool.Exec(ctx, "UPDATE user_deletions SET status = $1 WHERE id = $2", "scheduled", id)
			if err != nil {
				log.Printf("Error updating user_deletions status for UserID %s: %v", userID, err)
			}
			_ = incrementUserEmailCount(ctx, rc, userID)
			
		} else if status == "scheduled" && now.After(scheduledDate) {
			// Execute deletion if scheduled time has passed
			if !checkGlobalDeletionRateLimit(ctx, rc) {
				log.Printf("Skipping deletion for UserID %s: global rate limit exceeded", userID)
				continue
			}
			
			log.Printf("Executing deletion for UserID: %s. Scheduled for: %s", userID, scheduledDate.Format(time.RFC3339))
			
			// Execute soft delete
			err := executeUserDeletion(ctx, db, userID)
			if err != nil {
				log.Printf("Error executing deletion for UserID %s: %v", userID, err)
				continue
			}
			
			// Update status to executed
			_, err = db.Pool.Exec(ctx, "UPDATE user_deletions SET status = $1, executed = $2 WHERE id = $3", "executed", true, id)
			if err != nil {
				log.Printf("Error updating user_deletions status for UserID %s: %v", userID, err)
			}
			_ = incrementGlobalDeletionCount(ctx, rc)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating user_deletions rows: %v", err)
	}
}

// checkGlobalDeletionRateLimit checks if the global deletion rate limit has been exceeded.
func checkGlobalDeletionRateLimit(ctx context.Context, rc *redis.Client) bool {
	key := redisKeyGlobalDeletionCount + time.Now().Format("2006-01-02")
	count, err := rc.Get(ctx, key).Int64()
	if err != nil && err != redis.Nil {
		log.Printf("Error getting global deletion count from Redis: %v", err)
		return false // Fail safe: allow deletion if Redis is down or error
	}
	return count < maxGlobalDeletionsPerDay
}

// incrementGlobalDeletionCount increments the global deletion count.
func incrementGlobalDeletionCount(ctx context.Context, rc *redis.Client) error {
	key := redisKeyGlobalDeletionCount + time.Now().Format("2006-01-02")
	pipe := rc.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.ExpireAt(ctx, key, time.Now().Add(24*time.Hour).Truncate(24*time.Hour))
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Error incrementing global deletion count in Redis: %v", err)
		return err
	}
	log.Printf("Global deletion count incremented to %d", incr.Val())
	return nil
}

// checkUserEmailRateLimit checks if a user has exceeded their email rate limit.
func checkUserEmailRateLimit(ctx context.Context, rc *redis.Client, userID string) bool {
	key := redisKeyUserEmailCount + userID + ":" + time.Now().Format("2006-01-02")
	count, err := rc.Get(ctx, key).Int64()
	if err != nil && err != redis.Nil {
		log.Printf("Error getting user email count from Redis for UserID %s: %v", userID, err)
		return false // Fail safe: allow email if Redis is down or error
	}
	return count < maxEmailsPerUserPerDay
}

// incrementUserEmailCount increments the email count for a user.
func incrementUserEmailCount(ctx context.Context, rc *redis.Client, userID string) error {
	key := redisKeyUserEmailCount + userID + ":" + time.Now().Format("2006-01-02")
	pipe := rc.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.ExpireAt(ctx, key, time.Now().Add(24*time.Hour).Truncate(24*time.Hour))
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Error incrementing user email count in Redis for UserID %s: %v", userID, err)
		return err
	}
	log.Printf("User %s email count incremented to %d", userID, incr.Val())
	return nil
}

// executeUserDeletion performs the actual user deletion (soft delete)
func executeUserDeletion(ctx context.Context, db *infrastructure.DB, userID string) error {
	// Start a transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Soft delete the user
	_, err = tx.Exec(ctx, `
		UPDATE users 
		SET is_deleted = true, deleted_at = $1, deletion_due_at = $1
		WHERE id = $2 AND is_deleted = false
	`, time.Now(), userID)
	if err != nil {
		return err
	}

	// Log the deletion action
	_, err = tx.Exec(ctx, `
		INSERT INTO action_logs (user_id, action, meta, created_at)
		VALUES ($1, $2, $3, $4)
	`, userID, "user_deleted", `{"deleted_by": "system", "reason": "scheduled_deletion"}`, time.Now())
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit(ctx)
}
