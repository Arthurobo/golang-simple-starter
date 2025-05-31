package public

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", HomePage).Methods("GET")
	r.HandleFunc("/about", AboutPage).Methods("GET")
	r.HandleFunc("/health", HealthCheck).Methods("GET")
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "Welcome to the API"}`))
}

func AboutPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "About us page, could become an API, later"}`))
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "API is running"}`))
}
