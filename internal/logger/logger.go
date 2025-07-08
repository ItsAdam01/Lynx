package logger

import (
	"fmt"
	"log/slog"
	"os"
)

var globalLogger *slog.Logger

// InitLogger configures the global logger to write JSON to the specified file.
func InitLogger(logPath string) error {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	globalLogger = slog.New(handler)
	slog.SetDefault(globalLogger)

	return nil
}

// Info logs an informational message with structured key-value pairs.
func Info(msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.Info(msg, args...)
	} else {
		// Fallback if not initialized
		slog.Info(msg, args...)
	}
}

// Warn logs a warning message with structured key-value pairs.
func Warn(msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.Warn(msg, args...)
	} else {
		slog.Warn(msg, args...)
	}
}

// Error logs an error message with structured key-value pairs.
func Error(msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.Error(msg, args...)
	} else {
		slog.Error(msg, args...)
	}
}
