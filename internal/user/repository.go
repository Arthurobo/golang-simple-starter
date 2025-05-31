package user

import (
	"database/sql"
)

func GetAll(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetByID(db *sql.DB, id string) (User, error) {
	var u User
	err := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return u, err
	}
	return u, nil
}

func Create(db *sql.DB, u *User) error {
	return db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
}

func Update(db *sql.DB, id string, u *User) error {
	_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
	return err
}

func Delete(db *sql.DB, id string) (User, error) {
	u, err := GetByID(db, id)
	if err != nil {
		return User{}, err
	}
	_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
	return u, err
}
