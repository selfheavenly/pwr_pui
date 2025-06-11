package handlers

import (
	"PUI/serialization"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetTripDetails(c *gin.Context) {

	type TripDetailsGet struct {
		StopName    string `json:"stop_name"`
		RouteID     string `json:"tram_id"`
		DirectionID string `json:"direction_id"`
		ArrivalTime string `json:"arrival_time"`
	}

	var get TripDetailsGet
	if err := c.ShouldBindJSON(&get); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	//odds

	rates := []serialization.Odd{
		{Label: "0:00 - 0:30", Value: 1.78},
		{Label: "0:30 - 1:00", Value: 2.15},
		{Label: "1:00 - 1:30", Value: 1.95},
		{Label: "1:30 - 2:00", Value: 2.50},
	}

	//balance

	dbmpk, ok := c.MustGet("dbmpk").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	dbopen, ok := c.MustGet("dbopen").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	bimbom, exists := c.Get("google_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "google_id not found in context"})
		return
	}

	googleID, ok := bimbom.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid google_id type"})
		return
	}

	fmt.Println("gettripdetails googleid:", googleID)

	var balance float64
	err := dbmpk.QueryRow("SELECT balance FROM users WHERE google_id = ?", googleID).
		Scan(&balance)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var tripDetails serialization.TripDetails
	err = dbopen.QueryRow("SELECT DISTINCT st.stop_id, t.trip_headsign AS line\nFROM stop_times AS st\nJOIN trips AS t ON st.trip_id = t.trip_id\nJOIN routes AS r ON t.route_id = r.route_id\nJOIN stops AS s ON st.stop_id = s.stop_id\nWHERE\ns.stop_name = ?\nAND arrival_time = ?\nAND r.route_id = ? \nAND direction_id = ? \nAND t.service_id IN (\n        SELECT service_id FROM calendar\n        WHERE\n            (\n            (DAYOFWEEK(CURDATE()) = 1 AND sunday = 1) OR\n            (DAYOFWEEK(CURDATE()) = 2 AND monday = 1) OR\n            (DAYOFWEEK(CURDATE()) = 3 AND tuesday = 1) OR\n            (DAYOFWEEK(CURDATE()) = 4 AND wednesday = 1) OR\n            (DAYOFWEEK(CURDATE()) = 5 AND thursday = 1) OR\n            (DAYOFWEEK(CURDATE()) = 6 AND friday = 1) OR\n            (DAYOFWEEK(CURDATE()) = 7 AND saturday = 1)\n          )\n      )", get.StopName, get.ArrivalTime, get.RouteID, get.DirectionID).
		Scan(&tripDetails.StopID, &tripDetails.Line)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip Details not found"})
		return
	}

	tripDetails.TramID = get.RouteID
	tripDetails.StopName = get.StopName
	tripDetails.ArrivalTime = get.ArrivalTime
	tripDetails.Odds = rates
	tripDetails.Balance = balance

	c.JSON(http.StatusOK, tripDetails)
}

/*
type TripDetails struct {
TramID      string    `json:"route_id"` --
StopID      int       `json:"stop_id"` --
StopName    string    `json:"stop_name"` --
Line		string    `json:"line"` --
ArrivalTime string    `json:"arrival_time"` --
Odds        []Odd 	  `json:"odds"` --
Balance     int       `json:"balance"`
}
*/

func GetUserInfo(c *gin.Context) {
	db, ok := c.MustGet("dbmpk").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	bimbom, exists := c.Get("google_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "google_id not found in context"})
		return
	}

	googleID, ok := bimbom.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid google_id type"})
		return
	}

	var user serialization.Users
	err := db.QueryRow("SELECT user_id, google_id, email, name, balance FROM users WHERE google_id = ?", googleID).
		Scan(&user.UserID, &user.GoogleID, &user.Email, &user.Name, &user.Balance)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func Logout(c *gin.Context) {
	accessToken, exists := c.Get("access_token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not found"})
		return
	}

	revokeURL := "https://oauth2.googleapis.com/revoke?token=" + accessToken.(string)
	resp, err := http.Post(revokeURL, "application/x-www-form-urlencoded", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
		return
	}

	c.Set("google_id", nil)
	c.Set("access_token", nil)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func AddBalance(c *gin.Context) {
	db, ok := c.MustGet("dbmpk").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	bimbom, exists := c.Get("google_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "google_id not found in context"})
		return
	}

	googleID, ok := bimbom.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid google_id type"})
		return
	}

	var amount float64
	if err := c.ShouldBindJSON(&amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than zero"})
		return
	}

	_, err := db.Exec("UPDATE users SET balance = balance + ? WHERE google_id = ?", amount, googleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance updated successfully"})
}

// list of stops, each has a list of arriving trams, a name, (direction? possibly)
func GetStops(c *gin.Context) {
	db, ok := c.MustGet("dbopen").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	// Query to fetch stop details
	rows, err := db.Query("SELECT stop_id, stop_name FROM stops WHERE stop_id = 2863 OR stop_id = 1684 OR stop_id = 4899 OR stop_id = 1709")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stops"})
		return
	}
	defer rows.Close()

	var stopSummaries []serialization.StopSummary

	for rows.Next() {
		var stop serialization.StopSummary
		err := rows.Scan(&stop.StopID, &stop.StopName)
		if err != nil {
			fmt.Println("Error scanning stop details:", err)
			continue
		}

		// Query to fetch lines for the current stop
		lineRows, err := db.Query(`
			SELECT DISTINCT r.route_id
			FROM stop_times AS st
			JOIN trips AS t ON st.trip_id = t.trip_id
			JOIN routes AS r ON t.route_id = r.route_id
			WHERE st.stop_id = ?
			ORDER BY r.route_id`, stop.StopID)
		if err != nil {
			fmt.Println("Error fetching lines for stop:", err)
			continue
		}

		var lines []string
		for lineRows.Next() {
			var routeID string
			if err := lineRows.Scan(&routeID); err != nil {
				fmt.Println("Error scanning line:", err)
				continue
			}
			lines = append(lines, routeID)
		}
		lineRows.Close()

		stop.Lines = lines
		stopSummaries = append(stopSummaries, stop)
	}

	c.JSON(http.StatusOK, stopSummaries)
}

func GetStopInfo(c *gin.Context) {
	db, ok := c.MustGet("dbopen").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	stopID := c.Param("stopId")

	// 1. Get basic stop info
	var stop struct {
		Name string  `json:"name"`
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
	}
	err := db.QueryRow(`
		SELECT stop_name, stop_lat, stop_lon 
		FROM OpenWroclaw.stops 
		WHERE stop_id = ?`, stopID).Scan(&stop.Name, &stop.Lat, &stop.Lon)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
		return
	}

	// 2. Get all route_id + direction_id combos
	rows, err := db.Query(`
		SELECT DISTINCT r.route_id, direction_id
		FROM OpenWroclaw.stop_times AS st
		JOIN OpenWroclaw.trips AS t ON st.trip_id = t.trip_id
		JOIN OpenWroclaw.routes AS r ON t.route_id = r.route_id
		WHERE st.stop_id = ?
		ORDER BY r.route_id`, stopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch routes"})
		return
	}
	defer rows.Close()

	type RouteDir struct {
		RouteID      string   `json:"route_id"`
		DirectionID  int      `json:"direction_id"`
		ArrivalTimes []string `json:"arrival_times"`
	}
	var routeDirs []RouteDir

	// Get current weekday in SQL format
	weekday := time.Now().Weekday()
	weekdayColumn := map[time.Weekday]string{
		time.Sunday:    "sunday",
		time.Monday:    "monday",
		time.Tuesday:   "tuesday",
		time.Wednesday: "wednesday",
		time.Thursday:  "thursday",
		time.Friday:    "friday",
		time.Saturday:  "saturday",
	}[weekday]

	for rows.Next() {
		var routeID string
		var directionID int
		if err := rows.Scan(&routeID, &directionID); err != nil {
			continue
		}

		// 3. For each route+direction combo, get today's arrival_times
		arrivalQuery := fmt.Sprintf(`
			SELECT DISTINCT arrival_time
			FROM OpenWroclaw.stop_times AS st
			JOIN OpenWroclaw.trips AS t ON st.trip_id = t.trip_id
			JOIN OpenWroclaw.routes AS r ON t.route_id = r.route_id
			WHERE st.stop_id = ? AND r.route_id = ? AND direction_id = ? 
			AND t.service_id IN (
				SELECT service_id FROM OpenWroclaw.calendar WHERE %s = 1
			)
			ORDER BY arrival_time`, weekdayColumn)

		arrivalRows, err := db.Query(arrivalQuery, stopID, routeID, directionID)
		if err != nil {
			continue
		}

		var arrivalTimes []string
		for arrivalRows.Next() {
			var timeStr string
			arrivalRows.Scan(&timeStr)
			arrivalTimes = append(arrivalTimes, timeStr)
		}
		arrivalRows.Close()

		routeDirs = append(routeDirs, RouteDir{
			RouteID:      routeID,
			DirectionID:  directionID,
			ArrivalTimes: arrivalTimes,
		})
	}

	// Final JSON response
	c.JSON(http.StatusOK, gin.H{
		"stop":   stop,
		"routes": routeDirs,
	})
}

// all on bets for logged in user
func GetUserBets(c *gin.Context) {
	db, ok := c.MustGet("dbmpk").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	bimbom, exists := c.Get("google_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "google_id not found in context"})
		return
	}

	googleID, ok := bimbom.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid google_id type"})
		return
	}

	fmt.Println("getuserbets googleID:", googleID)

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// Fetch paginated bets
	rows, err := db.Query("SELECT bet_id, tram_line_id, stop_id, placed_at, bet_rate, bet_amount, status, predicted_time, google_id FROM bets WHERE google_id = ? LIMIT ? OFFSET ?", googleID, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch bets"})
		fmt.Println("Error fetching bets:", err)
		return
	}
	defer rows.Close()

	var bets []serialization.BetBrief
	for rows.Next() {
		var bet serialization.BetBrief
		err := rows.Scan(&bet.BetID, &bet.TramLineID, &bet.StopID, &bet.PlacedAt, &bet.BetRate, &bet.BetAmount, &bet.Status, &bet.PredictedTime, &bet.GoogleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get bets"})
			fmt.Println("Error scanning row:", err)
			return
		}
		bets = append(bets, bet)
	}

	// Get total count of bets
	var totalCount int
	err = db.QueryRow("SELECT COUNT(*) FROM bets WHERE google_id = ?", googleID).Scan(&totalCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch total count"})
		fmt.Println("Error fetching total count:", err)
		return
	}

	// Calculate total pages
	totalPages := totalCount / pageSize
	if totalCount%pageSize > 0 {
		totalPages++
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"bets":        bets,
		"page":        page,
		"page_size":   pageSize,
		"total":       totalCount,
		"total_pages": totalPages,
	})
}

/*

type BetBrief struct {
	BetID               int     `json:"bet_id"` --
	BetAmount           float64 `json:"bet_amount"` --
	BetRate             float64 `json:"bet_rate"` --
	PlacedAt            string  `json:"placed_at"` --
	Status              string  `json:"status"` --
	TramLineID          string  `json:"tram_lane_id"` --
	StopID              int     `json:"stop_id"` --
	StopName            string  `json:"stop_name"` mikolaj ma zrobic
}


*/

func PostBet(c *gin.Context) {
	db, ok := c.MustGet("dbmpk").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	var bet serialization.BetBrief

	if err := c.ShouldBindJSON(&bet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// tram line id -> route_id

	_, err := db.Query("INSERT INTO bets (tram_line_id, stop_id, placed_at, bet_rate, bet_amount, status, predicted_time, google_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", bet.TramLineID, bet.StopID, bet.PlacedAt, bet.BetRate, bet.BetAmount, bet.Status, bet.PredictedTime, bet.GoogleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place bet"})
		return
	}
}

type TramStopRates struct {
	StopID     string `json:"stop_id"`
	TramLaneID string `json:"tram_lane_id"`
}

func GetRates(c *gin.Context) {
	rates := []serialization.Odd{
		{Label: "0:00 - 0:30", Value: 1.78},
		{Label: "0:30 - 1:00", Value: 2.15},
		{Label: "1:00 - 1:30", Value: 1.95},
		{Label: "1:30 - 2:00", Value: 2.50},
	}

	c.JSON(http.StatusOK, rates)
}

/*
	export interface RateOdds {
	  label: string; // e.g. "0:00 - 0:30"
	  value: number; // e.g. 1.78
	}
*/

/*
	Bets

Won
Lost
Pending
PendingWin


*/

func validateBets(bets []serialization.BetBrief, c *gin.Context) {

	posMap := map[int]serialization.Position{
		1709: {Latitude: 51.11681942, Longitude: 17.05025444},
		1684: {Latitude: 51.10125726, Longitude: 17.10914151},
		4899: {Latitude: 51.11145165, Longitude: 17.06052919},
		2863: {Latitude: 51.09375561, Longitude: 16.98037615},
	}

	for _, bet := range bets {

		if bet.Status == "Won" {

		} else if bet.Status == "Lost" {

		} else if bet.Status == "Pending" {

			// get tram position
			positions, err := GetPosition(bet.TramLineID)
			if err != nil {
				fmt.Println("Error fetching tram positions:", err)
				continue
			}

			for _, position := range positions {
				if hasArrived(position, posMap[bet.StopID]) {
					bet.Status = "PendingWin"
				}
			}

		} else if bet.Status == "PendingWin" {

			calcReward(time.Now().String(), bet.PredictedTime, bet, c)

		} else {
			fmt.Println("Unknown bet status:", bet.Status)
		}

	}
}

func calcReward(arrivalTimeStr, predictedTimeStr string, bet serialization.BetBrief, c *gin.Context) (status string) {

	arrivalTime, _ := parseTime(arrivalTimeStr)
	predictedTime, _ := parseTime(predictedTimeStr)

	if (arrivalTime > predictedTime-1) || (arrivalTime < predictedTime+1) {
		status = "Won"

		amountWon := bet.BetAmount * bet.BetRate

		err := increaseBalance(bet.GoogleID, amountWon, c)
		if err != nil {
			return ""
		}
	} else {
		status = "Lost"
	}

	return
}

func increaseBalance(googleID string, amount float64, c *gin.Context) error {
	db, ok := c.MustGet("dbmpk").(*sql.DB)
	if !ok {
		return fmt.Errorf("failed to get database connection")
	}

	// Update the user's balance
	_, err := db.Exec("UPDATE users SET balance = balance + ? WHERE google_id = ?", amount, googleID)
	if err != nil {
		return fmt.Errorf("could not update user balance: %v", err)
	}

	return nil
}

func parseTime(timeStr string) (int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid delay format")
	}
	hours, err0 := strconv.Atoi(parts[0])
	mins, err1 := strconv.Atoi(parts[1])
	// secs, err2 := strconv.Atoi(parts[2])
	if err0 != nil || err1 != nil {
		return 0, fmt.Errorf("invalid delay numbers")
	}

	dS := hours*60 + mins*60

	return dS, nil
}

func getStopPosition(stopID string, c *gin.Context) (serialization.Position, error) {
	db, ok := c.MustGet("dbopen").(*sql.DB)
	if !ok {
		return serialization.Position{}, fmt.Errorf("failed to get database connection")
	}

	var position serialization.Position
	err := db.QueryRow("SELECT stop_lat, stop_lon FROM stops WHERE stop_id = ?", stopID).Scan(&position.Latitude, &position.Longitude)
	if err != nil {
		return serialization.Position{}, fmt.Errorf("could not fetch stop position: %v", err)
	}

	return position, nil
}

func parseDelayToSeconds(delayStr string) (int, error) {
	parts := strings.Split(delayStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid delay format")
	}
	mins, err1 := strconv.Atoi(parts[0])
	secs, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("invalid delay numbers")
	}

	dS := mins*60 + secs

	return dS, nil
}

/*
func groupBetsByInterval(bets []serialization.BetBrief, interval int) []serialization.Odd {
	binMap := make(map[int][]float64)

	for _, bet := range bets {
		seconds, err := parseDelayToSeconds(bet.ActualDelay)
		if err != nil {
			continue
		}
		bin := seconds / interval
		binMap[bin] = append(binMap[bin], bet.BetRate)
	}

	var result []serialization.Odd
	for bin, rates := range binMap {
		sum := 0.0
		for _, r := range rates {
			sum += r
		}
		avg := sum / float64(len(rates))
		label := fmt.Sprintf("%d:%02d - %d:%02d", (bin*interval)/60, (bin*interval)%60, ((bin+1)*interval)/60, ((bin+1)*interval)%60)
		result = append(result, serialization.Odd{
			Label: label,
			Value: avg,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Label < result[j].Label
	})

	return result
}
*/

func GetPosition(tramName string) (pos []serialization.VehiclePosition, err error) {
	apiURL := "https://mpk.wroc.pl/bus_position"
	data := url.Values{}
	data.Set("busList[tram][]", tramName)

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var positions []serialization.VehiclePosition
	if err := json.Unmarshal(body, &positions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	var tramPositions []serialization.VehiclePosition
	for _, pos := range positions {
		if pos.Name == tramName {
			tramPositions = append(tramPositions, pos)
		}
	}

	return tramPositions, nil
}

func hasArrived(tram serialization.VehiclePosition, currStop serialization.Position) bool {

	if tram.Latitude == 0 || tram.Longitude == 0 {
		return false
	}

	dist := calculateDistance(currStop.Latitude, currStop.Longitude, tram.Latitude, tram.Longitude)
	if dist < 0.0005 {
		return true
	}

	return false
}

// pythagorean theorem to calculate distance between two points on the earth
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {

	// Convert latitude and longitude differences to radians
	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)

	// Approximate distance using the Pythagorean theorem
	x := dLon * math.Cos((lat1+lat2)*(math.Pi/180)/2)
	y := dLat
	return math.Sqrt(x*x + y*y)
}

func checkTram(posTrams []serialization.VehiclePosition, startStopLoc, currStop, lastStopLoc serialization.Position) (bestTram serialization.VehiclePosition, err error) {

	type TramDistance struct {
		tram       serialization.VehiclePosition
		distToCurr float64
	}

	distanceToCurr := make([]TramDistance, 0)

	for _, tram := range posTrams {
		if tram.Latitude == 0 || tram.Longitude == 0 {
			continue
		}

		// calculate distance to the current stop
		dist := calculateDistance(currStop.Latitude, currStop.Longitude, tram.Latitude, tram.Longitude)

		distanceToCurr = append(distanceToCurr, TramDistance{tram: tram, distToCurr: dist})
	}

	// sort trams by distance to the current stop
	sort.Slice(TramDistance{}, func(i, j int) bool {
		return distanceToCurr[i].distToCurr < distanceToCurr[j].distToCurr
	})

	for _, tram := range distanceToCurr {
		startStopDist := calculateDistance(startStopLoc.Latitude, startStopLoc.Longitude, tram.tram.Latitude, tram.tram.Longitude)
		lastStopDist := calculateDistance(lastStopLoc.Latitude, lastStopLoc.Longitude, tram.tram.Latitude, tram.tram.Longitude)

		if startStopDist < lastStopDist {
			bestTram = tram.tram
			return bestTram, nil
		}
	}

	return bestTram, fmt.Errorf("no suitable tram found")
}

/*

func parseTime(timeStr string) (int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid delay format")
	}
	hours, err0 := strconv.Atoi(parts[0])
	mins, err1 := strconv.Atoi(parts[1])
	// secs, err2 := strconv.Atoi(parts[2])
	if err0 != nil || err1 != nil {
		return 0, fmt.Errorf("invalid delay numbers")
	}

	dS := hours*60 + mins*60

	return dS, nil
}

func calculateActualDelay(startStopLoc, currStopLoc serialization.Position, tram serialization.Position, routeID, directionID, stopID string, c *gin.Context) (delay int) {

	startStopDist := calculateDistance(startStopLoc.Latitude, startStopLoc.Longitude, tram.Latitude, tram.Longitude)
	currStopDist := calculateDistance(currStopLoc.Latitude, currStopLoc.Longitude, tram.Latitude, tram.Longitude)

	currTime, nextStopTime := GetNextStopTime(routeID, directionID, stopID, c)

	return 0
}

func GetNextStopTime(routeID, directionID, stopID string, c *gin.Context) (currTime, t int) {
	db, ok := c.MustGet("dbopen").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	rows, err := db.Query("SELECT DISTINCT arrival_time\nFROM stop_times AS st\nJOIN trips AS t ON st.trip_id = t.trip_id\nJOIN routes AS r ON t.route_id = r.route_id\nJOIN stops AS s ON st.stop_id = s.stop_id\nWHERE\ns.stop_id = ?\nAND r.route_id = ?\nAND direction_id = ?\nAND t.service_id IN (\n    SELECT service_id FROM calendar\n    WHERE\n        (\n        (DAYOFWEEK(CURDATE()) = 1 AND sunday = 1) OR\n        (DAYOFWEEK(CURDATE()) = 2 AND monday = 1) OR\n        (DAYOFWEEK(CURDATE()) = 3 AND tuesday = 1) OR\n        (DAYOFWEEK(CURDATE()) = 4 AND wednesday = 1) OR\n        (DAYOFWEEK(CURDATE()) = 5 AND thursday = 1) OR\n        (DAYOFWEEK(CURDATE()) = 6 AND friday = 1) OR\n        (DAYOFWEEK(CURDATE()) = 7 AND saturday = 1)\n      )\n  )\nORDER BY arrival_time;", stopID, routeID, directionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch bets"})
		return
	}
	defer rows.Close()

	var arrivalTimesString []string
	for rows.Next() {
		err := rows.Scan(&arrivalTimesString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get arrival times"})
			return
		}
	}

	var arrivalTime []int

	currTimeFake := time.Now()
	currTime = currTimeFake.Hour()*60 + currTimeFake.Minute()

	for _, at := range arrivalTimesString {
		arrival, err := parseTime(at)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse arrival time"})
			return
		}

		if arrival >= currTime {
			arrivalTime = append(arrivalTime, arrival)
		}
	}

	//find first time bigger than current time
	for _, at := range arrivalTime {
		if at >= currTime {
			t = at
			break
		}
	}

	return
}

*/
