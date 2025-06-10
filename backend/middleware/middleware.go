package middleware

import (
	"PUI/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware that checks the Authorization header and validates token with Google
func ValidateGoogleAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Malformed Authorization header"})
			return
		}

		// Validate token by requesting Google userinfo
		req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			return
		}

		var user models.User
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Google response"})
			return
		}

		// Store user info in context
		c.Set("google_id", user.ID)

		c.Next()
	}
}

func DatabaseMiddlewareMPK(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("dbmpk", db)
		c.Next()
	}
}

func DatabaseMiddlewareOpen(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("dbopen", db)
		c.Next()
	}
}
