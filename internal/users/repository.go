package users

import (
	"database/sql"
)

func GetAll(db *sql.DB) ([]GetAllUsersModel, error) {
	rows, err := db.Query("SELECT id, email, first_name, last_name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []GetAllUsersModel
	for rows.Next() {
		var u GetAllUsersModel
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetByID(db *sql.DB, id string) (IndividualUserModel, error) {
	var u IndividualUserModel
	err := db.QueryRow("SELECT id, email, first_name, last_name FROM users WHERE id = $1", id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName)
	if err != nil {
		return u, err
	}
	return u, nil
}

func GetUserByEmail(db *sql.DB, email string) (IndividualUserModel, error) {
	var u IndividualUserModel
	err := db.QueryRow("SELECT id, email, first_name, last_name, username FROM users WHERE email = $1", email).
		Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Username)

	if err == sql.ErrNoRows {
		return IndividualUserModel{}, nil // No user found, but this is not an error
	}
	return u, err // May still return a real error
}

func GetUserByUsername(db *sql.DB, username string) (IndividualUserModel, error) {
	var u IndividualUserModel
	err := db.QueryRow("SELECT id, email, first_name, last_name, username FROM users WHERE username = $1", username).
		Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Username)

	if err == sql.ErrNoRows {
		return IndividualUserModel{}, nil // No user found, but this is not an error
	}
	return u, err // May still return a real error
}

func Create(db *sql.DB, u *CreateUserModel) error {
	return db.QueryRow("INSERT INTO users (email, first_name, last_name, username) VALUES ($1, $2, $3, $4) RETURNING id", u.Email, u.FirstName, u.LastName, u.Username).Scan(&u.ID)
}

func Update(db *sql.DB, id string, u *UpdateUserModel) error {
	_, err := db.Exec("UPDATE users SET email = $1, first_name = $2, last_name = $3 WHERE id = $4", u.Email, u.FirstName, u.LastName, id)
	return err
}

func Delete(db *sql.DB, id string) (IndividualUserModel, error) {
	u, err := GetByID(db, id)
	if err != nil {
		return IndividualUserModel{}, err
	}
	_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
	return u, err
}
