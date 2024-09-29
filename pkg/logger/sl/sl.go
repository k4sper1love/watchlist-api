/*
Package sl provides logging functionality based on the `slog` package.
It supports various logging formats and levels tailored to different environments.
The package includes functions for setting up the logger and logging HTTP request details.
*/

package sl

import (
	"log/slog"
	"net/http"
	"os"
)

// Log is a global logger instance used throughout the application.
var Log *slog.Logger

// SetupLogger configures the global logger based on the specified environment.
// It sets different logging levels and formats for "local", "dev", and "prod" environments.
func SetupLogger(env string) {
	Log = configureLogger(env)
	Log = Log.With(slog.String("env", env))
}

// configureLogger initializes and returns a logger based on the provided environment.
func configureLogger(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case "local":
		// Logger for local environment with text output and debug level.
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "dev":
		// Logger for development environment with JSON output and debug level.
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "prod":
		// Logger for production environment with JSON output and info level.
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		// Default to production settings.
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	return slog.New(handler)
}

// PrintEndpointInfo logs information about an incoming HTTP request.
func PrintEndpointInfo(r *http.Request) {
	Log.Info(
		"handling request",
		slog.String("path", r.RequestURI),
		slog.String("method", r.Method),
		slog.String("from", r.RemoteAddr),
		slog.String("to", r.Host),
	)
}

// PrintEndpointError logs an error message related to an HTTP request.
func PrintEndpointError(msg string, err interface{}, r *http.Request) {
	Log.Error(
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
	Log.Warn(
		msg,
		slog.Any("error", err),
		slog.String("path", r.RequestURI),
		slog.String("method", r.Method),
		slog.String("from", r.RemoteAddr),
		slog.String("to", r.Host),
	)
}
