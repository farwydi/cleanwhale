package tonic

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/farwydi/cleanwhale/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	observedLogs, logs := observer.New(zap.InfoLevel)
	logger := zap.New(observedLogs)

	router.Use(NewRecoveryMiddleware(config.ModeLocal, logger))

	tests := map[string]interface{}{
		"str": "some panic",
		"err": errors.New("some error"),
	}
	for key, tc := range tests {
		tc := tc
		key := key
		t.Run(key, func(t *testing.T) {
			router.GET("/"+key, func(_ *gin.Context) {
				panic(tc)
			})

			req := httptest.NewRequest(http.MethodGet, "/"+key, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)

			require.Equal(t, 1, logs.Len())
			all := logs.TakeAll()

			logLine := all[0]
			assert.Equal(t, zap.ErrorLevel, logLine.Level)

			ctxMap := logLine.ContextMap()
			assert.NotZero(t, ctxMap["recover"])

			for _, field := range logLine.Context {
				if field.Key == "recover" {
					assert.Equal(t, tc, field.Interface)
				}
			}
		})
	}
}
