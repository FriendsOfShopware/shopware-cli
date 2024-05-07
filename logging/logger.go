package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

// loggerKey points to the value in the context where the logging is stored.
const loggerKey = contextKey("logging")

var fallbackLogger *zap.SugaredLogger

func NewLogger(verbose bool) *zap.SugaredLogger {
	loggerCfg := zap.NewDevelopmentConfig()
	loggerCfg.EncoderConfig.MessageKey = "message"
	loggerCfg.EncoderConfig.TimeKey = "timestamp"
	loggerCfg.EncoderConfig.EncodeDuration = zapcore.NanosDurationEncoder
	loggerCfg.EncoderConfig.StacktraceKey = "error.stack"
	loggerCfg.EncoderConfig.FunctionKey = "logging.method_name"
	loggerCfg.DisableStacktrace = true
	loggerCfg.DisableCaller = true
	loggerCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	if !verbose {
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		loggerCfg.EncoderConfig.TimeKey = ""
	}

	logger, err := loggerCfg.Build()
	if err != nil {
		logger = zap.NewNop()
	}

	return logger.Sugar()
}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}

	if fallbackLogger == nil {
		loggerCfg := zap.NewProductionConfig()
		logger, _ := loggerCfg.Build()

		fallbackLogger = logger.Sugar()
	}

	return fallbackLogger
}

func DisableLogger(ctx context.Context) context.Context {
	return WithLogger(ctx, zap.NewNop().Sugar())
}
