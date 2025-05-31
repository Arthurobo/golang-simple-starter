package user

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/mail"

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

func createHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		if _, err := mail.ParseAddress(u.Email); err != nil {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}
		if err := Create(db, &u); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	}
}

func updateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var u User
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
