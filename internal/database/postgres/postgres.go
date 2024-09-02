package postgres

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	_ "github.com/lib/pq"
	"log/slog"
)

var db *sql.DB

func OpenDB() (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", config.Dsn)
	if err != nil {
		sl.Log.Error("failed to open database connection", slog.Any("error", err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		sl.Log.Error("failed to ping database", slog.Any("error", err))
		return nil, err
	}

	sl.Log.Info("database connection established successfully")

	if config.Migrations != "" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			sl.Log.Error("failed to create migration driver", slog.Any("error", err))
			return nil, err
		}

		m, err := migrate.NewWithDatabaseInstance(config.Migrations, config.DB, driver)
		if err != nil {
			sl.Log.Error("failed to create migration instance", slog.Any("error", err))
			return nil, err
		}

		err = m.Up()
		if errors.Is(err, migrate.ErrNoChange) {
			sl.Log.Info("no migrations to apply")
		} else if err != nil {
			sl.Log.Error("migration failed", slog.Any("error", err))
			return nil, err
		} else {
			sl.Log.Info("migrations applied successfully")
		}
	}

	return db, nil
}
