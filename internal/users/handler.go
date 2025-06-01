package users

import (
	"api/pkg/middleware"
	"api/pkg/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/users/login", loginHandler(db)).Methods("POST")
	r.HandleFunc("/users", getAllUsersHandler(db)).Methods("GET")
	r.HandleFunc("/users", RegisterUserHandler(db)).Methods("POST")
	r.HandleFunc("/users/{id}", updateHandler(db)).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteHandler(db)).Methods("DELETE")

	// Protected routes
	r.Handle("/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(getUserByIDHandler(db)))).Methods("GET")
}

func getAllUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := GetAll(db)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error fetching users", err.Error())
			return
		}
		utils.WriteJSONSuccess(w, http.StatusOK, "Users fetched successfully", users)
	}
}

func getUserByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure token is valid and extract user ID
		userIDFromToken, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok {
			utils.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		// Optionally restrict access to only self
		idParam := mux.Vars(r)["id"]
		if fmt.Sprint(userIDFromToken) != idParam {
			utils.WriteJSONError(w, http.StatusForbidden, "Forbidden", nil)
			return
		}

		user, err := GetByID(db, idParam)
		if err != nil {
			utils.WriteJSONError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		utils.WriteJSONSuccess(w, http.StatusOK, "User fetched successfully", user)
	}
}

func RegisterUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u CreateUserModel
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request payload", nil)
			return
		}

		if errs := u.Validate(); len(errs) > 0 {
			utils.WriteJSONError(w, http.StatusBadRequest, "Validation failed", errs)
			return
		}

		existingUserEmail, err := GetUserByEmail(db, u.Email)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to check email", err.Error())
			return
		}
		if existingUserEmail.ID != 0 {
			utils.WriteJSONError(w, http.StatusBadRequest, "Email already exists", nil)
			return
		}

		existingUserUsername, err := GetUserByUsername(db, u.Username)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to check username", err.Error())
			return
		}
		if existingUserUsername.ID != 0 {
			utils.WriteJSONError(w, http.StatusBadRequest, "Username already exists", nil)
			return
		}

		if err := Create(db, &u); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to create user", err.Error())
			return
		}

		utils.WriteJSONSuccess(w, http.StatusCreated, "User Created successfully", u)
	}
}

func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequestModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request", nil)
			return
		}

		user, err := GetUserByEmail(db, req.Email)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
			utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid email or password", nil)
			return
		}

		accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email, user.FirstName, user.LastName)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Token generation failed", nil)
			return
		}

		utils.WriteJSONSuccess(w, http.StatusOK, "Login successful", map[string]interface{}{
			"access":     accessToken,
			"refresh":    refreshToken,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		})
	}
}

func updateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var u UpdateUserModel
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid payload", nil)
			return
		}
		if err := Update(db, id, &u); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to update user", err.Error())
			return
		}
		u.ID = id
		utils.WriteJSONSuccess(w, http.StatusOK, "User updated successfully", u)
	}
}

func deleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		_, err := Delete(db, id)
		if err != nil {
			utils.WriteJSONError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		utils.WriteJSONSuccess(w, http.StatusOK, "User deleted successfully", nil)
	}
}
