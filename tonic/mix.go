package tonic

import (
	"github.com/farwydi/cleanwhale/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewMix just gin and tonic.
// Creates a customized gin.Engine.
func NewMix(mode config.Mode, logger *zap.Logger) *gin.Engine {
	r := gin.New()

	if mode == config.ModeRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	AddRecoveryHandler(r, mode, logger)
	AddLoggerHandler(r, logger)
	AddMetricsHandler(r)
	AddHealthHandler(r)
	AddPProfHandler(r)

	return r
}
