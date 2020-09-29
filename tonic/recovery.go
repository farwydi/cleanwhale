package tonic

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/farwydi/cleanwhale/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// AddRecoveryHandler add to gin.Engine NewRecoveryMiddleware.
func AddRecoveryHandler(r *gin.Engine, mode config.Mode, logger *zap.Logger) {
	r.Use(NewRecoveryMiddleware(mode, logger))
}

// NewRecoveryMiddleware make handler recovered panic and logger him.
func NewRecoveryMiddleware(mode config.Mode, logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}

				mlog := logger.WithOptions(
					zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
				)
				if mode == config.ModeLocal {
					mlog = mlog.
						With(
							zap.Strings("headers", headers),
						)
				}

				mlog.Error("Recovery handler",
					zap.Reflect("recover", err))

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					c.Abort()
					return
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
