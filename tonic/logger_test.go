package tonic

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	observedLogs, logs := observer.New(zap.InfoLevel)
	logger := zap.New(observedLogs)

	router.Use(NewLoggerMiddleware(logger))

	router.GET("/", func(_ *gin.Context) {})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, 1, logs.Len())
	all := logs.TakeAll()

	logLine := all[0]
	assert.Equal(t, zap.InfoLevel, logLine.Level)

	ctxMap := logLine.ContextMap()
	assert.NotZero(t, ctxMap["client_ip"])
	assert.NotZero(t, ctxMap["method"])
	assert.NotZero(t, ctxMap["status_code"])
	assert.NotZero(t, ctxMap["body_size"])
	assert.NotZero(t, ctxMap["path"])

	for _, field := range logLine.Context {
		switch field.Key {
		case "client_ip":
			assert.Equal(t, "192.0.2.1", field.String)
		case "method":
			assert.Equal(t, http.MethodGet, field.String)
		case "status_code":
			assert.EqualValues(t, http.StatusOK, field.Integer)
		case "body_size":
			assert.EqualValues(t, -1, field.Integer)
		case "path":
			assert.Equal(t, "/", field.String)
		}
	}
}

func TestNewLoggerMiddlewareWithErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	observedLogs, logs := observer.New(zap.InfoLevel)
	logger := zap.New(observedLogs)

	router.Use(NewLoggerMiddleware(logger))

	err1 := errors.New("err 1")
	err2 := errors.New("err 2")

	router.GET("/", func(c *gin.Context) {
		_ = c.Error(err1)
		_ = c.Error(err2)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, 1, logs.Len())
	all := logs.TakeAll()

	logLine := all[0]
	assert.Equal(t, zap.ErrorLevel, logLine.Level)

	ctxMap := logLine.ContextMap()
	assert.NotZero(t, ctxMap["error"])
	assert.NotZero(t, ctxMap["client_ip"])
	assert.NotZero(t, ctxMap["method"])
	assert.NotZero(t, ctxMap["status_code"])
	assert.NotZero(t, ctxMap["body_size"])
	assert.NotZero(t, ctxMap["path"])

	for _, field := range logLine.Context {
		switch field.Key {
		case "error":
			assert.True(t, field.Equals(zap.Errors("error", []error{
				err1,
				err2,
			})))
		case "client_ip":
			assert.Equal(t, "192.0.2.1", field.String)
		case "method":
			assert.Equal(t, http.MethodGet, field.String)
		case "status_code":
			assert.EqualValues(t, http.StatusOK, field.Integer)
		case "body_size":
			assert.EqualValues(t, -1, field.Integer)
		case "path":
			assert.Equal(t, "/", field.String)
		}
	}
}
