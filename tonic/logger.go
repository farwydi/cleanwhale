package tonic

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AddLoggerHandler add to gin.Engine NewLoggerMiddleware.
func AddLoggerHandler(r *gin.Engine, logger *zap.Logger) {
	r.Use(NewLoggerMiddleware(logger))
}

// NewLoggerMiddleware make handler logs incoming requests.
func NewLoggerMiddleware(logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		requestURL := c.Request.URL.Path

		// Log only when path is not being skipped
		if inSkipHandler(requestURL) {
			c.Next()
			return
		}

		// Start timer
		start := time.Now()
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		if raw != "" {
			requestURL = requestURL + "?" + raw
		}

		logger := logger.With(
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.Int("status_code", c.Writer.Status()),
			zap.Int("body_size", c.Writer.Size()),
			zap.String("path", requestURL))

		errs := c.Errors.ByType(gin.ErrorTypePrivate)
		if len(errs) > 0 {
			logger.Error("Access",
				zap.Errors("error", mapperErrors(errs)))
			return
		}

		logger.Info("Access")
	}
}

func mapperErrors(errs []*gin.Error) (result []error) {
	result = make([]error, len(errs))
	for i, err := range errs {
		result[i] = err.Err
	}
	return result
}
