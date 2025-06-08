package handlers

import (
	"PUI/serialization"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserInfo(c *gin.Context) {

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
	err := db.QueryRow("SELECT stop_id, stop_name, stop_location_x, stop_location_y FROM stops WHERE stop_id = $1", stopID).
		Scan(&stop.StopID, &stop.StopName, &stop.StopLat, &stop.StopLon)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
		return
	}
}

// all on bets for logged in user
func GetUserBets(c *gin.Context) {
	db, ok := c.MustGet("db").(*sql.DB)
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
