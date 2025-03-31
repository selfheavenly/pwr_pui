package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" {
		log.Fatal("Missing database configuration: ensure DB_USER, DB_PASSWORD, DB_HOST, and DB_NAME are set")
	}

	// Connection string format: "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	defer db.Close()

	// Check connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Database is not reachable:", err)
	}

	fmt.Println("Successfully connected to MySQL!")

	selectQuery := "SELECT * FROM users"
	rows, err := db.Query(selectQuery)
	if err != nil {
		log.Fatal("Error querying data:", err)
	}

	fmt.Println("Fetched Users:")

	// Iterate over the rows and print the results
	for rows.Next() {
		var user_id, google_id int
		var email, name string
		var balance float32
		err := rows.Scan(&user_id, &google_id, &email, &name, &balance)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}

		fmt.Printf("user_id: %d, google_id: %d, email: %s, name %s, balance %f\n", user_id, google_id, email, name, balance)

	}

	// Check for errors from iteration
	if err := rows.Err(); err != nil {
		log.Fatal("Error in row iteration:", err)
	}

	defer rows.Close()
}
