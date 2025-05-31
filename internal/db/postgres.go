package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Connect(connStr string) (*sql.DB, error) {
	return sql.Open("postgres", connStr)
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL
		)
	`)
	return err
}
