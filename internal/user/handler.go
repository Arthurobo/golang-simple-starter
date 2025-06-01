package user

import (
	"api/pkg/utils"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/users", getAllHandler(db)).Methods("GET")
	r.HandleFunc("/users/{id}", getHandler(db)).Methods("GET")
	r.HandleFunc("/users", createHandler(db)).Methods("POST")
	r.HandleFunc("/users/{id}", updateHandler(db)).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteHandler(db)).Methods("DELETE")
}

func getAllHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := GetAll(db)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error fetching users", err.Error())
			return
		}
		utils.WriteJSONSuccess(w, http.StatusOK, "Users fetched successfully", users)
	}
}

func getHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		user, err := GetByID(db, id)
		if err != nil {
			utils.WriteJSONError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		utils.WriteJSONSuccess(w, http.StatusOK, "User fetched successfully", user)
	}
}

func createHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u CreateUser
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

func updateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var u UpdateUser
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
