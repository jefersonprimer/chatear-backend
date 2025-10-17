package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/supabase-community/supabase-go/supabase"
)


// NewRedisClient creates a new Redis client
func NewRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis URL: %w", err)
	}
	client := redis.NewClient(opt)
	// Ping to check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	return client, nil
}

// NewNatsClient creates a new NATS client
func NewNatsClient(natsURL string) (*nats.Conn, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}
	return conn, nil
}

// Infrastructure holds all infrastructure components
type Infrastructure struct {
	Supabase *supabase.Client
	Redis    *redis.Client
	NatsConn *nats.Conn
}

// NewInfrastructure creates and initializes all infrastructure components
func NewInfrastructure(supabaseURL, supabaseAnonKey, redisURL, natsURL string) (*Infrastructure, error) {
	var supabaseClient *supabase.Client
	var redisClient *redis.Client
	var natsConn *nats.Conn
	var err error

	supabaseClient, err = supabase.NewClient(supabaseURL, supabaseAnonKey)
	if err != nil {
		log.Printf("Failed to create Supabase client: %v", err)
		supabaseClient = nil
	}

	redisClient, err = NewRedisClient(redisURL)
	if err != nil {
		log.Printf("Failed to create redis client: %v", err)
		redisClient = nil
	}

	natsConn, err = NewNatsClient(natsURL)
	if err != nil {
		log.Printf("Failed to create nats client: %v", err)
		natsConn = nil
	}

	return &Infrastructure{
		Supabase: supabaseClient,
		Redis:    redisClient,
		NatsConn: natsConn,
	},
		nil
}

// Close closes all infrastructure connections
func (i *Infrastructure) Close() {
	if i.Redis != nil {
		i.Redis.Close()
	}
	if i.NatsConn != nil {
		i.NatsConn.Close()
	}
}
