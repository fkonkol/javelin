package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fkonkol/javelin/server/data"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

var (
	DB_URI string = os.Getenv("DB_URI")
)

func main() {
	log.Println("Server up and running")

	// Acquire connection with database
	pool := data.InitSQL(DB_URI)
	defer pool.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthCheck)

	log.Fatal(http.ListenAndServe(":8000", mux))
}
