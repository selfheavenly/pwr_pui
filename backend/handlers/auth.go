package handlers

import (
	"PUI/config"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"os"

	"github.com/gin-gonic/gin"
	//"golang.org/x/oauth2"
)

func HandleGoogleLogin(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL("randomstate")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Token exchange error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := config.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Println("Userinfo error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	db, ok := c.MustGet("db").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	err = AddUserToDB(db, userInfo["id"].(int), userInfo["id"].(string), userInfo["email"].(string), userInfo["name"].(string), 0.0)
	if err != nil {
		return
	}

	// Save to DB or session here...
	c.JSON(http.StatusOK, userInfo)
}

func AddUserToDB(db *sql.DB, userID int, googleID, email, name string, balance float64) error {
	// Check if the user exists
	query := "SELECT COUNT(*) FROM users WHERE user_id = ?"
	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}

	// If user does not exist, insert them into the database
	if count == 0 {
		insertQuery := "INSERT INTO users (user_id, google_id, email, name, balance) VALUES (?, ?, ?, ?, ?)"
		_, err := db.Exec(insertQuery, userID, googleID, email, name, balance)
		if err != nil {
			return fmt.Errorf("error adding user to database: %w", err)
		}
	}

	return nil
}
