// Package postgres provides functions for interacting with the PostgreSQL database.
// It includes operations for managing users, films, collections and permissions, and handling migrations.
//
// This package requires the `pq` and `golang-migrate` packages for PostgreSQL and migration support, respectively.
// Ensure that `config.Dsn` and `config.Migrations` are properly configured.

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

// OpenDB opens a PostgreSQL database connection and applies any pending migrations.
//
// It initializes a connection using the DSN from `config.Dsn`, verifies the connection with a ping,
// and applies migrations if a path is specified in `config.Migrations`.
func OpenDB() (*sql.DB, error) {
	var err error

	// Open a connection to the PostgreSQL database using the DSN from the configuration.
	db, err = sql.Open("postgres", config.Dsn)
	if err != nil {
		sl.Log.Error("failed to open database connection", slog.Any("error", err))
		return nil, err
	}

	// Ping the database to ensure the connection is valid.
	err = db.Ping()
	if err != nil {
		sl.Log.Error("failed to ping database", slog.Any("error", err))
		return nil, err
	}

	sl.Log.Info("database connection established successfully")

	// Check if migrations need to be applied.
	if config.Migrations != "" {
		// Create a migration driver instance using the database connection.
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			sl.Log.Error("failed to create migration driver", slog.Any("error", err))
			return nil, err
		}

		// Create a new migration instance with the specified migrations' path.
		m, err := migrate.NewWithDatabaseInstance(config.Migrations, "postgres", driver)
		if err != nil {
			sl.Log.Error("failed to create migration instance", slog.Any("error", err))
			return nil, err
		}

		// Apply any pending migrations.
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
