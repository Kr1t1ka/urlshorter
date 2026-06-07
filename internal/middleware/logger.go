package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		timeStart := time.Now()
		c.Next()
		duration := time.Since(timeStart)
		log.Info(
			"request",
			zap.String("URI", c.Request.RequestURI),
			zap.String("Method", c.Request.Method),
			zap.Duration("Duration", duration),
			zap.Int("status", c.Writer.Status()),
			zap.Int("Size", c.Writer.Size()),
		)
	}
}
