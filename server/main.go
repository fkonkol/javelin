package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	fmt.Println("Up and running")

	http.HandleFunc("/health", healthCheck)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
