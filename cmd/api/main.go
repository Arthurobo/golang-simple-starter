package main

import (
	"api/internal/db"
	"api/internal/posts"
	"api/internal/public"
	"api/internal/users"
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

	users.RegisterRoutes(router, database)
	public.RegisterRoutes(router)
	posts.RegisterRoutes(router, database)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
