package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fkonkol/javelin/server/account"
	"github.com/fkonkol/javelin/server/data"
	"github.com/fkonkol/javelin/server/messaging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

var (
	DB_URI = os.Getenv("DB_URI")
)

func main() {
	log.Println("Server up and running")

	// Acquire connection with database
	pool := data.InitSQL(DB_URI)
	defer pool.Close()

	// Redis cache used for storing account sessions
	sessions := data.InitSessionStore()

	// Initialize endpoint handlers
	accounts := account.NewHandler(pool, sessions)

	// Setup routes
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://localhost:3000"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/health", accounts.Auth(healthCheck))
	r.Post("/users/register", accounts.Register())
	r.Post("/users/login", accounts.Login())
	r.Get("/users", accounts.GetUserByUsername())
	r.Get("/auth", accounts.Auth(func(w http.ResponseWriter, r *http.Request) {}))
	r.Get("/ws", messaging.HandleNewConnection)

	log.Fatal(http.ListenAndServe(":8000", r))
}
