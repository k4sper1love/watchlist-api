// Package sl provides a structured logging system based on the `slog` package.
//
// Features:
// - Supports different logging formats (JSON, text) based on the environment.
// - Can log to a file, console, or both, depending on the `LOGS_OUTPUT` environment variable.
// - Uses `lumberjack` for log file rotation.
// - Provides helper functions for logging HTTP requests and errors.
//
// Configuration:
// - `LOGS_OUTPUT`: Defines where logs should be written (`file`, `console`, `both`).
// - `LOGS_DIR`: Specifies the directory where log files are stored.
// - `env`: Defines the logging environment (`local`, `dev`, `prod`).
package sl

import (
	"github.com/natefinch/lumberjack"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// Init initializes the logger based on the provided environment.
// It sets the global logger and configures output destinations.
func Init(env string) {
	handler := setupHandler(env)
	logger := slog.New(handler).With(slog.String("env", env))

	slog.SetDefault(logger)
}

// setupHandler configures and returns a logger handler based on the environment and log settings.
func setupHandler(env string) slog.Handler {
	logOutput := strings.ToLower(os.Getenv("LOGS_OUTPUT"))
	logsDir := os.Getenv("LOGS_DIR")

	if logsDir == "" {
		logsDir = "./logs" // Default log directory if not set
		slog.Warn("LOGS_DIR not set, using default", slog.String("path", logsDir))
	}

	logFile := &lumberjack.Logger{
		Filename:   logsDir + "/app.log",
		MaxSize:    10,   // Max file size in MB before rotation
		MaxBackups: 30,   // Max number of old log files
		Compress:   true, // Compress old logs
	}

	var output io.Writer
	switch logOutput {
	case "file":
		output = logFile
		slog.Info("Logging to file", slog.String("file", logFile.Filename))
	case "console":
		output = os.Stdout
		slog.Info("Logging to console")
	case "both":
		output = io.MultiWriter(os.Stdout, logFile)
		slog.Info("Logging to both console and file", slog.String("file", logFile.Filename))
	default:
		output = os.Stdout
		slog.Warn("Invalid LOGS_OUTPUT value, defaulting to console", slog.String("value", logOutput))
	}

	switch env {
	case "local":
		// Logger for local environment with text output and debug level.
		return slog.NewTextHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "dev":
		// Logger for development environment with JSON output and debug level.
		return slog.NewJSONHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "prod":
		// Logger for production environment with JSON output and info level.
		return slog.NewJSONHandler(output, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		// Default to local settings.
		slog.Warn("Unknown environment, defaulting to local settings", slog.String("env", env))
		return slog.NewTextHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
}

// PrintEndpointInfo logs information about an incoming HTTP request.
func PrintEndpointInfo(r *http.Request) {
	slog.Info(
		"handling request",
		slog.String("path", r.RequestURI),
		slog.String("method", r.Method),
		slog.String("from", r.RemoteAddr),
		slog.String("to", r.Host),
	)
}

// PrintEndpointError logs an error message related to an HTTP request.
func PrintEndpointError(msg string, err interface{}, r *http.Request) {
	slog.Error(
		msg,
		slog.Any("error", err),
		slog.String("path", r.RequestURI),
		slog.String("method", r.Method),
		slog.String("from", r.RemoteAddr),
		slog.String("to", r.Host),
	)
}

// PrintEndpointWarn logs a warning message related to an HTTP request.
func PrintEndpointWarn(msg string, err interface{}, r *http.Request) {
	slog.Warn(
		msg,
		slog.Any("error", err),
		slog.String("path", r.RequestURI),
		slog.String("method", r.Method),
		slog.String("from", r.RemoteAddr),
		slog.String("to", r.Host),
	)
}
