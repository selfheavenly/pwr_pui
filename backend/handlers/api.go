package handlers

import (
	"PUI/serialization"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

// get current logged in user info
func GetUserInfo(c *gin.Context) {
	db, ok := c.MustGet("db").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	// Call the serialization function
	userJSON, err := serialization.GetUsersJSON(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the serialized JSON
	c.JSON(http.StatusOK, gin.H{"data": userJSON})
}

// all info on one specific tram
func GetTramInfo(c *gin.Context) {
	db, ok := c.MustGet("db").(*sql.DB)
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
	db, ok := c.MustGet("db").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	rows, err := db.Query("SELECT stop_id, stop_name FROM stops_dictionary")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stops"})
		return
	}
	defer rows.Close()

	type StopWithTrams struct {
		Stop          serialization.StopsDictionary `json:"stop"`
		ArrivingTrams []string                      `json:"arriving_trams"`
	}

	var result []StopWithTrams

	for rows.Next() {
		var stop serialization.StopsDictionary
		err := rows.Scan(&stop.StopID, &stop.StopName)
		if err != nil {
			continue
		}

		tramRows, err := db.Query("SELECT tram_lane_id FROM stop_train_map WHERE stop_id = $1", stop.StopID)
		if err != nil {
			continue
		}

		var tramIDs []string
		for tramRows.Next() {
			var tramID string
			tramRows.Scan(&tramID)
			tramIDs = append(tramIDs, tramID)
		}
		tramRows.Close()

		result = append(result, StopWithTrams{
			Stop:          stop,
			ArrivingTrams: tramIDs,
		})
	}

	c.JSON(http.StatusOK, result)
}

// all about a specific stop
func GetStopInfo(c *gin.Context) {
	db, ok := c.MustGet("db").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	stopID := c.Param("id")

	var stop serialization.StopsDictionary
	err := db.QueryRow("SELECT stop_id, stop_name, stop_location_x, stop_location_y FROM stops_dictionary WHERE stop_id = $1", stopID).
		Scan(&stop.StopID, &stop.StopName, &stop.StopLocationX, &stop.StopLocationY)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
		return
	}

	rows, err := db.Query("SELECT tram_lane_id FROM stop_train_map WHERE stop_id = $1", stopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load trams for stop"})
		return
	}
	defer rows.Close()

	var trams []string
	for rows.Next() {
		var tramID string
		rows.Scan(&tramID)
		trams = append(trams, tramID)
	}

	c.JSON(http.StatusOK, gin.H{
		"stop":  stop,
		"trams": trams,
	})
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
