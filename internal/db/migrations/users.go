package migrations

import "database/sql"

func CreateUsersTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			username VARCHAR(255) NOT NULL UNIQUE,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			password TEXT NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
			date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}
