package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is global logger instance
var Log *zap.Logger

func init() {
	var err error

	switch os.Getenv("PHASE") {
	case "prod":
		config := zap.NewProductionConfig()
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.StacktraceKey = "" // error 발생 시, stack push를 비어있게함
		config.EncoderConfig = encoderConfig
		Log, err = config.Build(zap.AddCallerSkip(1))
	default:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		Log, err = config.Build()
	}

	if err != nil {
		panic(err)
	}
}

// Info is wrapper func for zap.Logger.Info
func Info(message string, fields ...zap.Field) {
	Log.Info(message, fields...)
}

// Debug is wrapper func for zap.Logger.Debug
func Debug(message string, fields ...zap.Field) {
	Log.Debug(message, fields...)
}

// Error is wrapper func for zap.Logger.Error
func Error(message string, fields ...zap.Field) {
	Log.Error(message, fields...)
}

/*
	// Logger usage from other files
	zap.L().Debug("this is hello func", zap.String("user", "leoo.j"), zap.Int("age", 28))
	zap.L().Info("this is hello func", zap.String("user", "leoo.j"), zap.Int("age", 28))
	zap.L().Error("this is hello func", zap.String("user", "leoo.j"), zap.Int("age", 28))

	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = "" // error 발생 시, stack push를 비어있게함
	config.EncoderConfig = encoderConfig

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	zap.L().Info("replaced zap's global loggers")
*/
