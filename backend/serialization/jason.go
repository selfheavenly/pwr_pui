package serialization

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	EntryID    int    `json:"entry_id"`
	TramLaneID string `json:"tram_lane_id"`
	StopID     int    `json:"stop_id,omitempty"`
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

type Calendar struct {
	ServiceID int  `json:"service_id"`
	Monday    *int `json:"monday,omitempty"`
	Tuesday   *int `json:"tuesday,omitempty"`
	Wednesday *int `json:"wednesday,omitempty"`
	Thursday  *int `json:"thursday,omitempty"`
	Friday    *int `json:"friday,omitempty"`
	Saturday  *int `json:"saturday,omitempty"`
	Sunday    *int `json:"sunday,omitempty"`
	StartDate *int `json:"start_date,omitempty"`
	EndDate   *int `json:"end_date,omitempty"`
}

// CalendarDates represents the calendar_dates table
type CalendarDates struct {
	ServiceID     int  `json:"service_id"`
	Date          *int `json:"date,omitempty"`
	ExceptionType *int `json:"exception_type,omitempty"`
}

// ControlStops represents the control_stops table
type ControlStops struct {
	VariantID *int `json:"variant_id,omitempty"`
	StopID    *int `json:"stop_id,omitempty"`
}

// Routes represents the routes table
type Routes struct {
	RouteID        *string `json:"route_id,omitempty"`
	AgencyID       *int    `json:"agency_id,omitempty"`
	RouteShortName *string `json:"route_short_name,omitempty"`
	RouteLongName  *string `json:"route_long_name,omitempty"`
	RouteDesc      *string `json:"route_desc,omitempty"`
	RouteType      *int    `json:"route_type,omitempty"`
	RouteType2ID   *int    `json:"route_type2_id,omitempty"`
	ValidFrom      *string `json:"valid_from,omitempty"`
	ValidUntil     *string `json:"valid_until,omitempty"`
}

// Shapes represents the shapes table
type Shapes struct {
	ShapeID         *int     `json:"shape_id,omitempty"`
	ShapePtLat      *float64 `json:"shape_pt_lat,omitempty"`
	ShapePtLon      *float64 `json:"shape_pt_lon,omitempty"`
	ShapePtSequence *int     `json:"shape_pt_sequence,omitempty"`
}

// StopTimes represents the stop_times table
type StopTimes struct {
	TripID        *string `json:"trip_id,omitempty"`
	ArrivalTime   *string `json:"arrival_time,omitempty"`
	DepartureTime *string `json:"departure_time,omitempty"`
	StopID        *int    `json:"stop_id,omitempty"`
	StopSequence  *int    `json:"stop_sequence,omitempty"`
	PickupType    *int    `json:"pickup_type,omitempty"`
	DropOffType   *int    `json:"drop_off_type,omitempty"`
}

// Stops represents the stops table
type Stops struct {
	StopID   *int     `json:"stop_id,omitempty"`
	StopCode *int     `json:"stop_code,omitempty"`
	StopName *string  `json:"stop_name,omitempty"`
	StopLat  *float64 `json:"stop_lat,omitempty"`
	StopLon  *float64 `json:"stop_lon,omitempty"`
}

// Trips represents the trips table
type Trips struct {
	RouteID      *string `json:"route_id,omitempty"`
	ServiceID    *int    `json:"service_id,omitempty"`
	TripID       *string `json:"trip_id,omitempty"`
	TripHeadsign *string `json:"trip_headsign,omitempty"`
	DirectionID  *int    `json:"direction_id,omitempty"`
	ShapeID      *int    `json:"shape_id,omitempty"`
	BrigadeID    *int    `json:"brigade_id,omitempty"`
	VehicleID    *int    `json:"vehicle_id,omitempty"`
	VariantID    *int    `json:"variant_id,omitempty"`
}

func GetUsersJSON(db *sql.DB) (string, error) {
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

func GetBetsJSON(db *sql.DB) (string, error) {
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

func GetStatusDictionaryJSON(db *sql.DB) (string, error) {
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

func GetStopTrainMapJSON(db *sql.DB) (string, error) {
	query := "SELECT entry_id, stop_id, tram_lane_id FROM stop_train_map"
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("error querying stop_train_map: %w", err)
	}
	defer rows.Close()

	var maps []StopTrainMap
	for rows.Next() {
		var stopMap StopTrainMap
		err := rows.Scan(&stopMap.StopID, &stopMap.StopID, &stopMap.TramLaneID)
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

func GetStopsDictionaryJSON(db *sql.DB) (string, error) {
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

func GetTramsDictionaryJSON(db *sql.DB) (string, error) {
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
