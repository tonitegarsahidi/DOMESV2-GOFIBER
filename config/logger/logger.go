package logger

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger(env string) {
	// Create logs directory if not exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Warning: Could not create logs directory: %v", err)
	}

	// Lumberjack daily/size logger configuration
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,  // megabytes per log file before rotating
		MaxBackups: 14,   // maximum number of old log files to retain
		MaxAge:     14,   // maximum number of days to retain old log files (14 days)
		Compress:   true, // compress rotated log files to gzip
		LocalTime:  true, // use local time for backups
	}

	// Write syncers
	fileWriter := zapcore.AddSync(lumberjackLogger)
	consoleWriter := zapcore.AddSync(os.Stdout)

	var consoleEncoder zapcore.Encoder
	var fileEncoder zapcore.Encoder
	var level zapcore.LevelEnabler

	if env == "production" {
		prodConfig := zap.NewProductionEncoderConfig()
		prodConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		
		consoleEncoder = zapcore.NewJSONEncoder(prodConfig)
		fileEncoder = zapcore.NewJSONEncoder(prodConfig)
		level = zap.InfoLevel
	} else {
		// Development console encoder (with colors)
		devConsoleConfig := zap.NewDevelopmentEncoderConfig()
		devConsoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		devConsoleConfig.TimeKey = "timestamp"
		devConsoleConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		consoleEncoder = zapcore.NewConsoleEncoder(devConsoleConfig)

		// Development file encoder (without colors to keep log files clean)
		devFileConfig := zap.NewDevelopmentEncoderConfig()
		devFileConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		devFileConfig.TimeKey = "timestamp"
		devFileConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		fileEncoder = zapcore.NewConsoleEncoder(devFileConfig)
		
		level = zap.DebugLevel
	}

	// Create unified core with multiple writers
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, level),
		zapcore.NewCore(fileEncoder, fileWriter, level),
	)

	// Build logger with caller info and automatic stacktrace for error-level logs
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
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
