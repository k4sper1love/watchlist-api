package sl

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var Log *slog.Logger

func SetupLogger(env string, file *os.File) {
	switch env {
	case "local":
		Log = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		Log = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		Log = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	Log = Log.With(slog.String("env", env))
}

func CreateLogFile(dir string) *os.File {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s/%s.log", dir, timestamp)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return file
}

func PrintHandlerInfo(r *http.Request) {
	Log.Info(
		"handling request",
		slog.String("path", r.RequestURI),
		slog.String("method", r.Method),
		slog.String("from", r.RemoteAddr),
		slog.String("to", r.Host),
	)

}
