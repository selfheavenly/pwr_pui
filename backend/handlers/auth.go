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
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Println("Error decoding user info:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	googleID, ok := userInfo["id"].(string)
	if !ok {
		log.Println("Invalid google_id in user info")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set("google_id", googleID)

	// Add user to the database if not already present
	db, ok := c.MustGet("dbmpk").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
		return
	}

	email := userInfo["email"].(string)
	name := userInfo["name"].(string)
	balance := 0.0

	if err := AddUserToDB(db, googleID, email, name, balance); err != nil {
		log.Println("Error adding user to database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User authenticated successfully", "user": userInfo})
}

func AddUserToDB(db *sql.DB, googleID, email, name string, balance float64) error {
	// Check if the user exists
	query := "SELECT COUNT(*) FROM users WHERE google_id = ?"
	var count int
	err := db.QueryRow(query, googleID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}

	// If user does not exist, insert them into the database
	if count == 0 {
		insertQuery := "INSERT INTO users (google_id, email, name, balance) VALUES (?, ?, ?, ?)"
		_, err := db.Exec(insertQuery, googleID, email, name, balance)
		if err != nil {
			return fmt.Errorf("error adding user to database: %w", err)
		}
	}

	return nil
}
