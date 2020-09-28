package log

import (
	"github.com/farwydi/cleanwhale/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap.Logger. The config.ProjectConfig variable sets the
// base logger params like project name and config.Mode.
// The CoreMakeFunc enumeration custom zapcore.Core for the logger.
func NewLogger(project config.ProjectConfig, cores ...CoreMakeFunc) (*zap.Logger, error) {
	zapCores := make([]zapcore.Core, 0, len(cores)+1)

	if len(cores) == 0 {
		switch project.Mode {
		case config.ModeLocal:
			cores = append(cores, DefaultDevelopmentCore)
		case config.ModeRelease:
			cores = append(cores, DefaultProductionCore)
		case config.ModeTest:
			return zap.NewNop(), nil
		default:
			cores = append(cores, DefaultDevelopmentCore)
		}
	}

	for _, coreFunc := range cores {
		core, err := coreFunc(project.Name)
		if err != nil {
			return nil, err
		}

		zapCores = append(zapCores, core)
	}

	logger := zap.New(
		zapcore.NewTee(zapCores...),
	)

	return logger, nil
}

// RegisterLogger logger for access by zap.L()
// and redirects stdlog to zap by default.
func RegisterLogger(logger *zap.Logger) {
	_ = zap.ReplaceGlobals(logger)
	_ = zap.RedirectStdLog(logger)
}
