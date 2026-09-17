package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if !validRequestID(requestID) {
			requestID = uuid.NewString()
		}
		c.Set("request_id", requestID)
		c.Header(requestIDHeader, requestID)
		c.Next()
	}
}

func validRequestID(requestID string) bool {
	if len(requestID) > 128 {
		return false
	}
	_, err := uuid.Parse(requestID)
	return err == nil
}
