/*
Package sl provides logging functionality based on the `slog` package.
It supports various logging formats and levels tailored to different environments.
The package includes functions for setting up the logger, creating log files, and logging HTTP request details.
*/

package sl

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Log is a global logger instance used throughout the application.
var Log *slog.Logger

// SetupLogger configures the global logger based on the environment and log file.
// It sets different logging levels and formats for "local", "dev", and "prod" environments.
func SetupLogger(env string, file *os.File) {
	switch env {
	case "local":
		// Configure logger for local environment with text output and debug level.
		Log = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		// Configure logger for development environment with JSON output and debug level.
		Log = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		// Configure logger for production environment with JSON output and info level.
		Log = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	// Add environment tag to the logger for additional context.
	Log = Log.With(slog.String("env", env))
}

// CreateLogFile creates and opens a log file in the specified directory with a timestamp in the filename.
// It returns the file handle for further use.
func CreateLogFile(dir string) *os.File {
	// Generate a timestamped filename for the log file.
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s/%s.log", dir, timestamp)

	// Open or create the log file with append mode and read/write permissions.
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Panic if the file cannot be opened or created.
		panic(err)
	}

	// Return the file handle for the created log file.
	return file
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
