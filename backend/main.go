package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Users represents the users table
type Users struct {
	UserID   int      `json:"user_id"`
	GoogleID *string  `json:"google_id,omitempty"`
	Email    *string  `json:"email,omitempty"`
	Name     *string  `json:"name,omitempty"`
	Balance  *float64 `json:"balance,omitempty"`
}

// Bets represents the bets table
type Bets struct {
	BetID         int     `json:"bet_id"`
	UserID        int     `json:"user_id"`
	TramLaneID    string  `json:"tram_lane_id"`
	StopID        int     `json:"stop_id"`
	BetDelay      *string `json:"bet_delay,omitempty"`
	BetTime       *string `json:"bet_time,omitempty"`
	BetMultiplier *string `json:"bet_multiplier,omitempty"`
	BetStatus     string  `json:"bet_status"`
}

// StatusDictionary represents the status_dictionary table
type StatusDictionary struct {
	Status string `json:"status"`
}

// StopTrainMap represents the stop_train_map table
type StopTrainMap struct {
	StopID     int     `json:"stop_id"`
	TramLaneID string  `json:"tram_lane_id"`
	LastSeen   *string `json:"last_seen,omitempty"`
}

// StopsDictionary represents the stops_dictionary table
type StopsDictionary struct {
	StopID        int     `json:"stop_id"`
	StopLocationX *string `json:"stop_location_x,omitempty"`
	StopLocationY *string `json:"stop_location_y,omitempty"`
	StopName      *string `json:"stop_name,omitempty"`
}

// TramsDictionary represents the trams_dictionary table
type TramsDictionary struct {
	TramLaneID string  `json:"tram_lane_id"`
	TramName   *string `json:"tram_name,omitempty"`
}

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

	jsonUsers, err := getUsersJSON(db)
	if err != nil {
		log.Fatal(err)
	}

	jsonBets, err := getBetsJSON(db)
	if err != nil {
		log.Fatal(err)
	}

	jsonStatusDictionary, err := getStatusDictionaryJSON(db)
	if err != nil {
		log.Fatal(err)
	}

	jsonStopStrainMap, err := getStopTrainMapJSON(db)
	if err != nil {
		log.Fatal(err)
	}

	jsonStopsDictionary, err := getStopsDictionaryJSON(db)
	if err != nil {
		log.Fatal(err)
	}

	jsonTramsDictionary, err := getTramsDictionaryJSON(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(jsonUsers)
	fmt.Println(jsonBets)
	fmt.Println(jsonStatusDictionary)
	fmt.Println(jsonStopStrainMap)
	fmt.Println(jsonStopsDictionary)
	fmt.Println(jsonTramsDictionary)

}

func getUsersJSON(db *sql.DB) (string, error) {
	selectQuery := "SELECT user_id, google_id, email, name, balance FROM users"
	rows, err := db.Query(selectQuery)
	if err != nil {
		return "", fmt.Errorf("error querying data: %w", err)
	}
	defer rows.Close()

	var users []Users

	// Iterate over the rows and populate the users slice
	for rows.Next() {
		var user Users
		err := rows.Scan(&user.UserID, &user.GoogleID, &user.Email, &user.Name, &user.Balance)
		if err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}
		users = append(users, user)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating rows: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting to JSON: %w", err)
	}

	return string(jsonData), nil
}

func getBetsJSON(db *sql.DB) (string, error) {
	query := "SELECT bet_id, user_id, tram_lane_id, stop_id, bet_delay, bet_time, bet_multiplier, bet_status FROM bets"
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("error querying bets: %w", err)
	}
	defer rows.Close()

	var bets []Bets
	for rows.Next() {
		var bet Bets
		err := rows.Scan(&bet.BetID, &bet.UserID, &bet.TramLaneID, &bet.StopID, &bet.BetDelay, &bet.BetTime, &bet.BetMultiplier, &bet.BetStatus)
		if err != nil {
			return "", fmt.Errorf("error scanning bet: %w", err)
		}
		bets = append(bets, bet)
	}

	jsonData, err := json.MarshalIndent(bets, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting bets to JSON: %w", err)
	}

	return string(jsonData), nil
}

func getStatusDictionaryJSON(db *sql.DB) (string, error) {
	query := "SELECT status FROM status_dictionary"
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("error querying status dictionary: %w", err)
	}
	defer rows.Close()

	var statuses []StatusDictionary
	for rows.Next() {
		var status StatusDictionary
		err := rows.Scan(&status.Status)
		if err != nil {
			return "", fmt.Errorf("error scanning status: %w", err)
		}
		statuses = append(statuses, status)
	}

	jsonData, err := json.MarshalIndent(statuses, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting status dictionary to JSON: %w", err)
	}

	return string(jsonData), nil
}

func getStopTrainMapJSON(db *sql.DB) (string, error) {
	query := "SELECT stop_id, tram_lane_id, last_seen FROM stop_train_map"
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("error querying stop_train_map: %w", err)
	}
	defer rows.Close()

	var maps []StopTrainMap
	for rows.Next() {
		var stopMap StopTrainMap
		err := rows.Scan(&stopMap.StopID, &stopMap.TramLaneID, &stopMap.LastSeen)
		if err != nil {
			return "", fmt.Errorf("error scanning stop_train_map: %w", err)
		}
		maps = append(maps, stopMap)
	}

	jsonData, err := json.MarshalIndent(maps, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting stop_train_map to JSON: %w", err)
	}

	return string(jsonData), nil
}

func getStopsDictionaryJSON(db *sql.DB) (string, error) {
	query := "SELECT stop_id, stop_location_x, stop_location_y, stop_name FROM stops_dictionary"
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("error querying stops_dictionary: %w", err)
	}
	defer rows.Close()

	var stops []StopsDictionary
	for rows.Next() {
		var stop StopsDictionary
		err := rows.Scan(&stop.StopID, &stop.StopLocationX, &stop.StopLocationY, &stop.StopName)
		if err != nil {
			return "", fmt.Errorf("error scanning stops_dictionary: %w", err)
		}
		stops = append(stops, stop)
	}

	jsonData, err := json.MarshalIndent(stops, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting stops_dictionary to JSON: %w", err)
	}

	return string(jsonData), nil
}

func getTramsDictionaryJSON(db *sql.DB) (string, error) {
	query := "SELECT tram_lane_id, tram_name FROM trams_dictionary"
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("error querying trams_dictionary: %w", err)
	}
	defer rows.Close()

	var trams []TramsDictionary
	for rows.Next() {
		var tram TramsDictionary
		err := rows.Scan(&tram.TramLaneID, &tram.TramName)
		if err != nil {
			return "", fmt.Errorf("error scanning trams_dictionary: %w", err)
		}
		trams = append(trams, tram)
	}

	jsonData, err := json.MarshalIndent(trams, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting trams_dictionary to JSON: %w", err)
	}

	return string(jsonData), nil
}
