package watchlist

import (
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/transport/rest"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"log/slog"
	"os"
)

func Run(args []string) {
	sl.SetupLogger("local", os.Stdout)
	sl.Log.Info("starting application")

	err := config.ParseEnv()
	if err != nil {
		sl.Log.Error("failed to load .env file", slog.Any("error", err))
		os.Exit(1)
	}

	sl.Log.Debug("environment variables parsed successfully")

	err = config.ParseFlags(args[1:])
	if err != nil {
		sl.Log.Error("failed to load flags", slog.Any("error", err))
		os.Exit(1)
	}

	sl.Log.Debug("command-line flags parsed successfully")

	sl.SetupLogger(config.Env, os.Stdout)

	db, err := postgres.OpenDB()
	if err != nil {
		os.Exit(1)
	}

	defer func() {
		if err := db.Close(); err != nil {
			sl.Log.Error("failed to close database connection", slog.Any("error", err))
			os.Exit(1)
		}
		sl.Log.Info("database connection closed")
	}()

	err = rest.Serve()
	if err != nil {
		os.Exit(1)
	}
}
