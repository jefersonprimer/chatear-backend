package middleware

import (
	"github.com/gin-gonic/gin"
	customhttp "github.com/jefersonprimer/chatear-backend/presentation/http"
)

// GinContextToContextMiddleware adds the Gin context to the request context.
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := customhttp.WithGinContext(c.Request.Context(), c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
