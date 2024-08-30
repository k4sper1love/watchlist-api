package postgres

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/k4sper1love/watchlist-api/internal/config"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

func OpenDB() (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", config.Dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if config.Migrations != "" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, err
		}

		m, err := migrate.NewWithDatabaseInstance(config.Migrations, config.DB, driver)
		if err != nil {
			return nil, err
		}

		err = m.Up()
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")
		} else if err != nil {
			log.Printf("migration failed: %v", err)
			return nil, err
		} else {
			log.Println("migrations applied successfully")
		}
	}

	return db, nil
}
