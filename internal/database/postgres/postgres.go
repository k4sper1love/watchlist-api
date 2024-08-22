package postgres

import (
	"database/sql"
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	_ "github.com/lib/pq"
)

func connectPostgres() *sql.DB {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Hostname, config.Port, config.Username, config.Password, config.Database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil
	}

	return db
}
