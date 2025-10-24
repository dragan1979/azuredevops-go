package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB is the globally accessible database connection handle.
var DB *sql.DB

// User struct to map data from the database.
// This matches your init.sql schema including created_at and is_active.
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

// InitDB initializes the database connection with a retry loop for resilience.
func InitDB() error {
	// Retrieve connection parameters from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return fmt.Errorf("missing one or more required database environment variables")
	}

	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	log.Printf("Attempting to connect to: %s@tcp(%s:%s)/%s", dbUser, dbHost, dbPort, dbName)

	// Open the database connection (this does not verify connectivity yet)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool limits
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// === RESILIENCY: Connection Retry Loop ===
	const maxRetries = 12 // 12 retries * 5 seconds = 60 seconds max wait
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err == nil {
			// Success! Assign the active connection and exit the loop.
			DB = db
			log.Println("Successfully connected to MySQL!")
			return nil
		}
		log.Printf("Database ping failed (Attempt %d/%d). Retrying in 5 seconds...", i+1, maxRetries)
		time.Sleep(5 * time.Second)
	}
	// If the loop finishes without success, return the final error.
	return fmt.Errorf("failed to establish database connection after %d retries", maxRetries)
}

// GetUsers retrieves all users from the 'users' table.
func GetUsers() ([]User, error) {
	// Ensure the query matches the Go User struct fields
	rows, err := DB.Query("SELECT id, username, email, created_at, is_active FROM users")
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		// Ensure the scan order matches the SELECT column order
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.IsActive); err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return users, nil
}
