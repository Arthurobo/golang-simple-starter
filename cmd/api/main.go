package main

import (
	"api/internal/db"
	"api/internal/public"
	"api/internal/user"
	"api/pkg/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	database, err := db.Connect(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.Use(middleware.JSONMiddleware)

	user.RegisterRoutes(router, database)
	public.RegisterRoutes(router)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
