package log

import "go.uber.org/zap/zapcore"

type (
	// CoreMakeFunc maker function type for zapcore.Core.
	CoreMakeFunc func(projectName string) (zapcore.Core, error)
)
