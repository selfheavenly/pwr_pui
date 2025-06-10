package handlers

import (
	"PUI/serialization"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetUserInfo(c *gin.Context) {
	db, ok := c.MustGet("dbmpk").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	googleID_get, exists := c.Get("google_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "google_id not found in context"})
		return
	}

	googleID, ok := googleID_get.(string)
	if !ok || googleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid google_id"})
		return
	}

	err := db.Ping()
	if err != nil {
		log.Fatal("Database is not reachable:", err)
	}

	var user serialization.Users
	err = db.QueryRow("SELECT user_id, google_id, email, name, balance FROM users WHERE google_id = ?", googleID).
		Scan(&user.UserID, &user.GoogleID, &user.Email, &user.Name, &user.Balance)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// all info on one specific tram
func GetTramInfo(c *gin.Context) {
	db, ok := c.MustGet("dbmpk").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	tramID := c.Param("id")
	var tram serialization.TramsDictionary

	err := db.QueryRow("SELECT tram_lane_id, tram_name FROM trams_dictionary WHERE tram_lane_id = $1", tramID).
		Scan(&tram.TramLaneID, &tram.TramName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tram not found"})
		return
	}

	c.JSON(http.StatusOK, tram)
}

// list of stops, each has a list of arriving trams, a name, (direction? possibly)
func GetStops(c *gin.Context) {
	db, ok := c.MustGet("dbopen").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	rows, err := db.Query("SELECT stop_id, stop_name, stop_lat, stop_lon FROM stops")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stops"})
		return
	}
	defer rows.Close()

	var stops []serialization.Stops
	for rows.Next() {
		var stop serialization.Stops
		if err := rows.Scan(&stop.StopID, &stop.StopName, &stop.StopLat, &stop.StopLon); err == nil {
			stops = append(stops, stop)
		}
	}

	c.JSON(http.StatusOK, stops)
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

	userID := c.GetInt("user_id") // Assumes middleware sets this

	rows, err := db.Query("SELECT bet_id, user_id, tram_lane_id, stop_id, bet_delay, bet_time, bet_multiplier, bet_status FROM bets WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch bets"})
		return
	}
	defer rows.Close()

	var bets []serialization.Bets

	for rows.Next() {
		var bet serialization.Bets
		err := rows.Scan(&bet.BetID, &bet.UserID, &bet.TramLaneID, &bet.StopID, &bet.BetDelay, &bet.BetTime, &bet.BetMultiplier, &bet.BetStatus)
		if err == nil {
			bets = append(bets, bet)
		}
	}

	c.JSON(http.StatusOK, bets)
}

func GetBetInfo(c *gin.Context) {

}

func PostBet(c *gin.Context) {

}

func GetRates(c *gin.Context) {

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
