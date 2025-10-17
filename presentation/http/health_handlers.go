package http

import (
	"context"
	"net"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jefersonprimer/chatear-backend/config"
	"github.com/jefersonprimer/chatear-backend/infrastructure"
	"github.com/nats-io/nats.go"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	infra *infrastructure.Infrastructure
	cfg   *config.Config
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(infra *infrastructure.Infrastructure, cfg *config.Config) *HealthHandler {
	return &HealthHandler{infra: infra, cfg: cfg}
}

// Healthz is a liveness probe to check if the service is running.
func (h *HealthHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}

// Readyz is a readiness probe to check if the service is ready to handle requests.
func (h *HealthHandler) Readyz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	errors := make(map[string]string)

	// Check Postgres
	if h.infra.Postgres != nil {
		if err := h.infra.Postgres.Health(ctx); err != nil {
			errors["postgres"] = err.Error()
		}
	} else {
		errors["postgres"] = "not configured or failed to connect"
	}

	// Check Redis
	if h.infra.Redis != nil {
		if err := h.infra.Redis.Ping(ctx).Err(); err != nil {
			errors["redis"] = err.Error()
		}
	} else {
		errors["redis"] = "not configured or failed to connect"
	}

	// Check NATS
	if h.infra.NatsConn != nil {
		if h.infra.NatsConn.Status() != nats.CONNECTED {
			errors["nats"] = "not connected"
		}
	} else {
		errors["nats"] = "not configured or failed to connect"
	}

	// Check SMTP
	if h.cfg.SMTPHost != "" {
		addr := h.cfg.SMTPHost + ":" + strconv.Itoa(h.cfg.SMTPPort)
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err != nil {
			errors["smtp"] = "failed to connect: " + err.Error()
		} else {
			client, err := smtp.NewClient(conn, h.cfg.SMTPHost)
			if err != nil {
				errors["smtp"] = "failed to create client: " + err.Error()
			} else {
				if h.cfg.SMTPUser != "" && h.cfg.SMTPPass != "" {
					auth := smtp.PlainAuth("", h.cfg.SMTPUser, h.cfg.SMTPPass, h.cfg.SMTPHost)
					if err := client.Auth(auth); err != nil {
						errors["smtp"] = "authentication failed: " + err.Error()
					}
				}
				client.Close()
			}
		}
	} else {
		errors["smtp"] = "not configured"
	}

	if len(errors) > 0 {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
			"errors": errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}