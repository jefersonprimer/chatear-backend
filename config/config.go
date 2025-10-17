package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all the configuration for the application
type Config struct {
	AppURL                  string
	Port                    int
	SupabaseConnectionString             string
	RedisURL                string
	NatsURL                 string
	JwtSecret               string
	AccessTokenTTL          time.Duration
	RefreshTokenTTL         time.Duration
	SMTPHost                string
	SMTPPort                int
	SMTPUser                string
	SMTPPass                string
	SMTPFrom                string
	MagicLinkExpiry         time.Duration
	RateLimitEnabled        bool
	KeyRotationInterval     time.Duration
	MaxEmailsPerDay         int
	HardDeleteRetentionPeriod time.Duration
}

// LoadConfig loads the configuration from the environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		AppURL:                  getEnv("APP_URL"),
		Port:                    getEnvAsInt("PORT"),
		SupabaseConnectionString:             getEnv("SUPABASE_CONNECTION_STRING", ""),
		RedisURL:                getEnv("REDIS_URL"),
		NatsURL:                 getEnv("NATS_URL""),
		JwtSecret:               getEnv("JWT_SECRET"),
		AccessTokenTTL:          getEnvAsDuration("ACCESS_TOKEN_TTL"),
		RefreshTokenTTL:         getEnvAsDuration("REFRESH_TOKEN_TTL"),
		SMTPHost:                getEnv("SMTP_HOST"),
		SMTPPort:                getEnvAsInt("SMTP_PORT"),
		SMTPUser:                getEnv("SMTP_USER"),
		SMTPPass:                getEnv("SMTP_PASS"),
		SMTPFrom:                getEnv("SMTP_FROM"),
		MagicLinkExpiry:         getEnvAsDuration("MAGIC_LINK_EXPIRY"),
		RateLimitEnabled:        getEnvAsBool("RATE_LIMIT_ENABLED"),
		KeyRotationInterval:     getEnvAsDuration("KEY_ROTATION_INTERVAL"),
		MaxEmailsPerDay:         getEnvAsInt("MAX_EMAILS_PER_DAY"),
		HardDeleteRetentionPeriod: getEnvAsDuration("HARD_DELETE_RETENTION_PERIOD"),
	}
}

// Helper functions to get environment variables

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}
