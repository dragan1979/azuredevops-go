package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	// Must be imported for `main` to use the functions defined in db.go
	// and to ensure the structs are available.
)

// Define a struct for the data passed to the HTML template
type TemplateData struct {
	Users []User
}

// GetRootHandler handles the root path ("/") and renders the styled HTML template.
func GetRootHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch data from the database using the function in db.go
	users, err := GetUsers()
	if err != nil {
		log.Printf("Error fetching users for template: %v", err)
		http.Error(w, "Failed to load user data for web page.", http.StatusInternalServerError)
		return
	}

	// 2. Parse the HTML template (CRITICAL: Parsing inside the handler ensures changes are picked up on every request for development)
	// We use Must(ParseGlob) here for simplicity in parsing.
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		// Return a user-friendly error if the template file can't be read/parsed
		http.Error(w, "Could not load application template. Check templates/index.html.", http.StatusInternalServerError)
		return
	}

	// 3. Render the template with the fetched data
	data := TemplateData{Users: users}
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render page.", http.StatusInternalServerError)
	}
}

// GetUsersHandler handles the /users API endpoint and returns raw JSON data.
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header for JSON response
	w.Header().Set("Content-Type", "application/json")

	users, err := GetUsers()
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		// Respond with a JSON error message
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve users"})
		return
	}

	// Respond with the list of users as JSON
	json.NewEncoder(w).Encode(users)
}