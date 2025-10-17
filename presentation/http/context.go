package http

import (
	"context"

	"github.com/gin-gonic/gin"
)

type contextKey string

const ginContextKey = contextKey("GinContext")

// WithGinContext returns a new context with the gin.Context stored in it.
func WithGinContext(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ginContextKey, c)
}

// GinContextFromContext returns the gin.Context from the context.
func GinContextFromContext(ctx context.Context) (*gin.Context, bool) {
	c, ok := ctx.Value(ginContextKey).(*gin.Context)
	return c, ok
}
