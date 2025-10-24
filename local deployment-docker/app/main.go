package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// 1. Initialize the Database Connection.
	if err := InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Ensure the database connection is closed when main exits
	defer DB.Close() 

	// 2. Setup Router
	r := mux.NewRouter()
	
	// NEW: Define the handler for the root path ("/") to serve the HTML template
	r.HandleFunc("/", GetRootHandler).Methods("GET")

	// Define the handler for fetching users (the JSON API endpoint)
	r.HandleFunc("/users", GetUsersHandler).Methods("GET")
	
	// Define a simple health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is healthy and connected to DB."))
	}).Methods("GET")

	// 3. Start Server
	port := "8080"
	log.Printf("Server listening on http://0.0.0.0:%s", port)
	
	// Start serving HTTP requests
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
