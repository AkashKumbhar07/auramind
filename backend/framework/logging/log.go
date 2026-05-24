package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string, format string) *zap.Logger {
	var lvl zapcore.Level
	switch level {
	case "debug":
		lvl = zapcore.DebugLevel
	case "info":
		lvl = zapcore.InfoLevel
	case "warn":
		lvl = zapcore.WarnLevel
	case "error":
		lvl = zapcore.ErrorLevel
	default:
		lvl = zapcore.InfoLevel
	}

	var enc zapcore.Encoder
	switch format {
	case "json":
		enc = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	default:
		enc = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), lvl)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func NewProduction() *zap.Logger {
	return New("info", "json")
}

func NewDevelopment() *zap.Logger {
	return New("debug", "console")
}
