package middleware

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

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
		defer resp.Body.Close()

		var user GoogleUser
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Google response"})
			return
		}

		// Store user info in context
		c.Set("google_user", user)
		c.Set("user_email", user.Email)
		c.Set("user_id", user.ID)

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

func GoogleIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		googleID := session.Get("google_id")

		if googleID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Set google_id in the context
		c.Set("google_id", googleID)
		c.Next()
	}
}
