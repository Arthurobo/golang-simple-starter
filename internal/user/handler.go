package user

import (
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
			http.Error(w, "Error fetching users", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	}
}

func getHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		user, err := GetByID(db, id)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string, details interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"error":   details,
	})
}

func createHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u CreateUser
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid request payload", nil)
			return
		}

		if errs := u.Validate(); len(errs) > 0 {
			writeJSONError(w, http.StatusBadRequest, "Validation failed", errs)
			return
		}

		existingUserEmail, err := GetUserByEmail(db, u.Email)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to check email", err.Error())
			return
		}
		if existingUserEmail.ID != 0 {
			writeJSONError(w, http.StatusBadRequest, "Email already exists", nil)
			return
		}

		existingUserUsername, err := GetUserByUsername(db, u.Username)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to check username", err.Error())
			return
		}
		if existingUserUsername.ID != 0 {
			writeJSONError(w, http.StatusBadRequest, "Username already exists", nil)
			return
		}

		if err := Create(db, &u); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to create user", err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	}
}

func updateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var u UpdateUser
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}
		if err := Update(db, id, &u); err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(u)
	}
}

func deleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		u, err := Delete(db, id)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User deleted successfully",
			"user":    u,
		})
	}
}
