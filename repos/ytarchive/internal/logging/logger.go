package logging

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	// Default to JSON handler for production
	logFormat := os.Getenv("LOG_FORMAT")
	logLevel := os.Getenv("LOG_LEVEL")

	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	if logFormat == "text" {
		Logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	} else {
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	slog.SetDefault(Logger)
}

// Convenience functions
func Info(msg string, args ...any)  { Logger.Info(msg, args...) }
func Error(msg string, args ...any) { Logger.Error(msg, args...) }
func Debug(msg string, args ...any) { Logger.Debug(msg, args...) }
func Warn(msg string, args ...any)  { Logger.Warn(msg, args...) }

// WithComponent returns a logger with a component attribute
func WithComponent(component string) *slog.Logger {
	return Logger.With("component", component)
}
