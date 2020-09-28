package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DefaultProductionCore for production config.Mode.
// By default json output format and iso8601 time.
// Minimal level at zap.InfoLevel.
var DefaultProductionCore CoreMakeFunc = func(_ string) (zapcore.Core, error) {
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	return core, nil
}

// DefaultDevelopmentCore for development config.Mode.
// By default console output format with color and iso8601 time.
// Minimal level at zap.DebugLevel.
var DefaultDevelopmentCore CoreMakeFunc = func(projectName string) (zapcore.Core, error) {
	encoder := zap.NewDevelopmentEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	// paints levels in colors
	encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	return core, nil
}
