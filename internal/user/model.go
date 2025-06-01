package user

import "api/pkg/validators"

type CreateUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (u *CreateUser) Validate() []string {
	var errors []string
	if u.Email == "" {
		errors = append(errors, "email is required")
	}
	if !validators.IsValidEmail(u.Email) {
		errors = append(errors, "invalid email")
	}
	if u.Username == "" {
		errors = append(errors, "username is required")
	}
	if u.FirstName == "" {
		errors = append(errors, "first_name is required")
	}
	if u.LastName == "" {
		errors = append(errors, "last_name is required")
	}
	return errors
}

type UpdateUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GetAllUsers struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type IndividualUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}
