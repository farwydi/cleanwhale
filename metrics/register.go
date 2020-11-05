// Package metrics helper functions for prometheus and basic metrics.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	isSetup = false

	// RequestCount counts the total number of requests.
	RequestCount *prometheus.CounterVec

	// RequestDuration counts the execution time of requests.
	RequestDuration *prometheus.SummaryVec

	// ResponseSize counts the number of bytes sent.
	ResponseSize *prometheus.SummaryVec

	// RequestSize counts the number of bytes received.
	RequestSize *prometheus.SummaryVec
)

// IsSetupCalled Returns the status of the method RegisterMetrics call.
func IsSetupCalled() bool {
	return isSetup
}

// RegisterMetrics base http metrics in prometheus.
func RegisterMetrics(subsystem string, logger *zap.Logger) {
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "requests_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "version", "method", "host", "handler"},
	)
	if err := prometheus.Register(RequestCount); err != nil {
		logger.Fatal("fail register metrics",
			zap.String("metric", "requests_total"),
			zap.Error(err))
	}

	RequestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Subsystem: subsystem,
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
		},
		[]string{"code", "version", "method", "host", "handler"},
	)
	if err := prometheus.Register(RequestDuration); err != nil {
		logger.Fatal("fail register metrics",
			zap.String("metric", "request_duration_seconds"),
			zap.Error(err))
	}

	ResponseSize = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Subsystem: subsystem,
			Name:      "response_size_bytes",
			Help:      "The HTTP response sizes in bytes.",
		},
		[]string{"code", "version", "method", "host", "handler"},
	)
	if err := prometheus.Register(ResponseSize); err != nil {
		logger.Fatal("fail register metrics",
			zap.String("metric", "response_size_bytes"),
			zap.Error(err))
	}

	RequestSize = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Subsystem: subsystem,
			Name:      "request_size_bytes",
			Help:      "The HTTP request sizes in bytes.",
		},
		[]string{"code", "version", "method", "host", "handler"},
	)
	if err := prometheus.Register(RequestSize); err != nil {
		logger.Fatal("fail register metrics",
			zap.String("metric", "request_size_bytes"),
			zap.Error(err))
	}

	isSetup = true
}
