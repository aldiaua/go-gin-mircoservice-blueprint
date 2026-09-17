package httpmiddleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AccessLog(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()
		log.WithFields(logrus.Fields{
			"request_id":  c.GetString("request_id"),
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status":      c.Writer.Status(),
			"duration_ms": time.Since(started).Milliseconds(),
		}).Info("http request")
	}
}
