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

// SetupLogger configures the global logger based on the environment.
// It sets different logging levels and formats for "local", "dev", and "prod" environments.
func SetupLogger(env string) {
	switch env {
	case "local":
		// Configure logger for local environment with text output and debug level.
		Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		// Configure logger for development environment with JSON output and debug level.
		Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		// Configure logger for production environment with JSON output and info level.
		Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	// Add environment tag to the logger for additional context.
	Log = Log.With(slog.String("env", env))
}

// PrintHandlerInfo logs details about an incoming HTTP request.
// It includes the request URI, method, remote address, and host.
func PrintHandlerInfo(r *http.Request) {
	Log.Info(
		"handling request",
		slog.String("path", r.RequestURI),
		slog.String("method", r.Method),
		slog.String("from", r.RemoteAddr),
		slog.String("to", r.Host),
	)
}

func PrintHandlerError(msg string, err interface{}, r *http.Request) {
	Log.Error(
		msg,
		slog.Any("error", err),
		slog.String("path", r.RequestURI),
		slog.String("method", r.Method),
		slog.String("from", r.RemoteAddr),
		slog.String("to", r.Host),
	)
}

func PrintHandlerWarn(msg string, err interface{}, r *http.Request) {
	Log.Warn(
		msg,
		slog.Any("error", err),
		slog.String("path", r.RequestURI),
		slog.String("method", r.Method),
		slog.String("from", r.RemoteAddr),
		slog.String("to", r.Host),
	)
}
