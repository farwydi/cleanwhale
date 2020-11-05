package tonic

import (
	"net/http"
	"strconv"
	"time"

	"github.com/farwydi/cleanwhale/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const metricsPath = "/metrics"

// AddMetricsHandler add to gin.Engine NewMetricsMiddleware.
func AddMetricsHandler(r *gin.Engine) {
	r.Use(NewMetricsMiddleware())
	r.GET(metricsPath, gin.WrapH(promhttp.Handler()))
}

// NewMetricsMiddleware make handler with registered basic metrics.
func NewMetricsMiddleware() func(c *gin.Context) {
	if metrics.IsSetupCalled() {
		zap.L().Warn("Metrics is not register!, Call metrics.RegisterMetrics(...)")
	}

	return func(c *gin.Context) {
		requestURL := c.Request.URL.String()

		// Log only when path is not being skipped
		if inSkipHandler(requestURL) {
			c.Next()
			return
		}

		requestSize := float64(ComputeApproximateRequestSize(c.Request))

		// Call next handler
		timeBeforeStartNextHandler := time.Now()
		c.Next()

		secondsSinceNext := time.Since(timeBeforeStartNextHandler).Seconds()

		responseStatus := strconv.Itoa(c.Writer.Status())

		responseSize := float64(c.Writer.Size())

		requestURL = c.GetString("handler.url")
		if requestURL == "" {
			// Protected memory leak
			requestURL = c.HandlerName()
		}

		version := c.GetString("handler.version")
		if len(version) == 0 {
			version = "vX"
		}

		labels := prometheus.Labels{
			"code":    responseStatus,
			"version": version,
			"method":  c.Request.Method,
			"host":    c.Request.Host,
			"handler": requestURL,
		}

		if metrics.RequestCount != nil {
			metrics.RequestCount.With(labels).Inc()
		}

		if metrics.RequestDuration != nil {
			metrics.RequestDuration.With(labels).Observe(secondsSinceNext)
		}

		if metrics.RequestSize != nil {
			metrics.RequestSize.With(labels).Observe(requestSize)
		}

		if metrics.ResponseSize != nil {
			metrics.ResponseSize.With(labels).Observe(responseSize)
		}
	}
}

// ComputeApproximateRequestSize determines the actual number of bytes transferred.
// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func ComputeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
