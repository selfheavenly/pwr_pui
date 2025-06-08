package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func Connect() (db *sql.DB) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME_MPK")

	fmt.Printf("DB_USER: %s, DB_PASSWORD: %s, DB_HOST: %s, DB_NAME_OPEN: %s\n", dbUser, dbPassword, dbHost, dbName)

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" {
		log.Fatal("Missing database configuration: ensure DB_USER, DB_PASSWORD, DB_HOST, and DB_NAME are set")
	}

	// Connection string format: "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
	fmt.Println("Connection string:", dsn)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Check connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Database is not reachable:", err)
	}

	fmt.Println("Successfully connected to MPK MySQL!")

	return
}

func ConnectOpen() (db *sql.DB) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME_OPEN")

	// Log environment variables for debugging
	fmt.Printf("DB_USER: %s, DB_PASSWORD: %s, DB_HOST: %s, DB_NAME_OPEN: %s\n", dbUser, dbPassword, dbHost, dbName)

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" {
		log.Fatal("Missing database configuration: ensure DB_USER, DB_PASSWORD, DB_HOST, and DB_NAME_OPEN are set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
	fmt.Println("Connection string:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Database is not reachable:", err)
	}

	fmt.Println("Successfully connected to Open MySQL!")

	return
}
