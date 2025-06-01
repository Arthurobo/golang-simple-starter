package users

import "api/pkg/validators"

type LoginRequestModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type CreateUserModel struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
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
	if u.Password == "" {
		errors = append(errors, "password is required")
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
	Username  string `json:"username"`
	IsActive  bool   `json:"is_active"`
}

type IndividualUserModel struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"-"` // This is not returned in the response, using it for comparison
	IsActive  bool   `json:"is_active"`
	IsDeleted bool   `json:"-"`
}
