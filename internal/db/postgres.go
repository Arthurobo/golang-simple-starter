package db

import (
	"api/internal/db/migrations"
	"database/sql"

	_ "github.com/lib/pq"
)

func Connect(connStr string) (*sql.DB, error) {
	return sql.Open("postgres", connStr)
}

func Migrate(db *sql.DB) error {
	if err := migrations.CreateUsersTable(db); err != nil {
		return err
	}
	if err := migrations.CreatePostsTable(db); err != nil {
		return err
	}
	return nil
}
