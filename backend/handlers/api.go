package handlers

import (
	"PUI/serialization"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

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

// list of stops, each has a list of arriving trams, a name, (direction? possibly)
func GetStops(c *gin.Context) {
	db, ok := c.MustGet("dbopen").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	rows, err := db.Query("")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stops"})
		return
	}
	defer rows.Close()

	var StopSummary []serialization.StopSummary

	for rows.Next() {
		var stop serialization.StopSummary
		var linesRaw string
		// Adjust the query and scan as needed for your schema
		err := rows.Scan(&stop.StopID, &stop.StopName, &linesRaw)
		if err == nil {
			// Assuming linesRaw is a comma-separated string of line names/numbers
			stop.Lines = strings.Split(linesRaw, ",")
			StopSummary = append(StopSummary, stop)
		}
	}

	c.JSON(http.StatusOK, StopSummary)
}

// all about a specific stop
func GetStopInfo(c *gin.Context) {
	db, ok := c.MustGet("db").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	stopID := c.Param("id")

	var stop serialization.Stops
	err := db.QueryRow("SELECT DISTINCT\n"+
		"    r.route_short_name,"+
		"    s.stop_name,"+
		"    s.stop_lat,"+
		"    s.stop_lon"+
		"FROM"+
		"    stop_times AS st"+
		"JOIN"+
		"    trips AS t ON st.trip_id = t.trip_id"+
		"JOIN"+
		"    routes AS r ON t.route_id = r.route_id"+
		"JOIN"+
		"    stops AS s ON st.stop_id = s.stop_id"+
		"WHERE s.stop_id = $1"+
		"ORDER BY"+
		"    route_short_name;", stopID).
		Scan(&stop.StopID, &stop.StopName, &stop.StopLat, &stop.StopLon)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
		return
	}
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

	rows, err := db.Query("SELECT bet_id, tram_lane_id, stop_id, bet_delay, bet_time, bet_multiplier, bet_status FROM bets WHERE user_id = ?", googleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch bets"})
		return
	}
	defer rows.Close()

	var bets []serialization.RecentBetsBrief

	for rows.Next() {
		var bet serialization.BetBrief
		err := rows.Scan(&bet.BetID, &bet.TramLaneID, &bet.StopID, &bet.BetAmount, &bet.PlacedAt, &bet.BetRate, &bet.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get bets"})
			return
		}
	}

	c.JSON(http.StatusOK, bets)
}

/*

type BetBrief struct {
	BetID               int     `json:"bet_id"` --
	BetAmount           float64 `json:"bet_amount"` --
	BetRate             float64 `json:"bet_rate"` --
	PlacedAt            string  `json:"placed_at"` --
	Status              string  `json:"status"` --
	TramLaneID          string  `json:"tram_lane_id"` --
	StopID              int     `json:"stop_id"` --
	StopName            string  `json:"stop_name"` mikolaj ma zrobic
	// TODO
	ActualDelay         string  `json:"actual_delay"`
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

	_, err := db.Query("INSERT INTO bets (tram_lane_id, stop_id, bet_delay, bet_time, bet_multiplier, bet_status, stop_name) VALUES (?, ?, ?, ?, ?, ?, ?)", bet.TramLaneID, bet.StopID, bet.BetAmount, bet.PlacedAt, bet.BetRate, bet.Status, bet.StopName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place bet"})
		return
	}
}

func GetRates(c *gin.Context) {

}

/*
	export interface RateOdds {
	  label: string; // e.g. "0:00 - 0:30"
	  value: number; // e.g. 1.78
	}

chuj
*/
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

/*
SELECT DISTINCT
    s.stop_id,
    r.route_short_name,
    st.arrival_time,
    st.departure_time
FROM
    stop_times AS st
JOIN
    trips AS t ON st.trip_id = t.trip_id
JOIN
    routes AS r ON t.route_id = r.route_id
JOIN
    stops AS s ON st.stop_id = s.stop_id
WHERE s.stop_id = "15" AND r.route_short_name = "123"
ORDER BY
    arrival_time;
*/
