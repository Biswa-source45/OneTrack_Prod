package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuditMiddleware captures request context for audit logging
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate unique request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Capture client IP
		c.Set("client_ip", c.ClientIP())

		// Capture user agent
		c.Set("user_agent", c.Request.UserAgent())

		// Continue processing
		c.Next()
	}
}
