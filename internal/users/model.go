package users

import "api/pkg/validators"

type CreateUserModel struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (u *CreateUserModel) Validate() []string {
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

type UpdateUserModel struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GetAllUsersModel struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type IndividualUserModel struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}
