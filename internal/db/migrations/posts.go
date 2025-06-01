package migrations

import "database/sql"

func CreatePostsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			is_deleted BOOLEAN DEFAULT FALSE,
			date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}
