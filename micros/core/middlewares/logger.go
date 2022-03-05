package middlewares

import (
	customLogger "bridge/service-managers/logger"
	"io"
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func httpLogger(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {
	logger := customLogger.Get().
		With().
		Int("status", c.Writer.Status()).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("ip", c.ClientIP()).
		Dur("latency", latency).
		Str("user_agent", c.Request.UserAgent()).
		Logger()
	return logger
}

func LoggerMw() gin.HandlerFunc {
	return logger.SetLogger(logger.WithLogger(httpLogger))
}
