package logger

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(env string) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.OutputPaths = []string{"logs/app.log", "stdout"}
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	}

	// Create logs directory if not exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Warning: Could not create logs directory: %v", err)
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	zap.ReplaceGlobals(Logger)
}

func GetLogger() *zap.Logger {
	return Logger
}

func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Convenience functions
func Info(msg string, fields ...zapcore.Field) {
	Logger.Info(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	Logger.Fatal(msg, fields...)
}

func Debug(msg string, fields ...zapcore.Field) {
	Logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	Logger.Warn(msg, fields...)
}
