package middleware

import (
	"github.com/gin-gonic/gin"

	"go-restful/lib/utils"
)

// RequestID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get(utils.XRequestID)

		// Initial request id with UUID
		if requestID == "" {
			requestID = utils.GenRequestID()
		}

		// Expose it for use in the application
		c.Set(utils.XRequestID, requestID)

		// Set X-Request-ID header
		c.Writer.Header().Set(utils.XRequestID, requestID)
		c.Next()
	}
}
