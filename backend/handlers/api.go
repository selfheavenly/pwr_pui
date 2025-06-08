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

func GetTramInfo(c *gin.Context) {

}

func GetStops(c *gin.Context) {

}

func GetStopInfo(c *gin.Context) {

}

func GetUserBets(c *gin.Context) {

}

func GetBetInfo(c *gin.Context) {

}

func PostBet(c *gin.Context) {

}

func GetRates(c *gin.Context) {

}
