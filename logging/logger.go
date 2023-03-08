package logging

import (
	"context"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

// loggerKey points to the value in the context where the logging is stored.
const loggerKey = contextKey("logging")

var fallbackLogger *zap.SugaredLogger

func NewLogger(debug bool) *zap.SugaredLogger {
	loggerCfg := zap.NewDevelopmentConfig()
	loggerCfg.EncoderConfig.EncodeTime = timeEncoder()
	loggerCfg.EncoderConfig.MessageKey = "message"
	loggerCfg.EncoderConfig.TimeKey = "timestamp"
	loggerCfg.EncoderConfig.EncodeDuration = zapcore.NanosDurationEncoder
	loggerCfg.EncoderConfig.StacktraceKey = "error.stack"
	loggerCfg.EncoderConfig.FunctionKey = "logging.method_name"
	loggerCfg.DisableStacktrace = !debug

	if debug {
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := loggerCfg.Build()
	if err != nil {
		logger = zap.NewNop()
	}

	return logger.Sugar()
}

// timeEncoder encodes the time as RFC3339 nano.
func timeEncoder() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339Nano))
	}
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
