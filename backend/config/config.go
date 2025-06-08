package config

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig *oauth2.Config

func LoadEnv() {

	/*
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

	*/

	GoogleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	if GoogleOauthConfig.ClientID == "" || GoogleOauthConfig.ClientSecret == "" || GoogleOauthConfig.RedirectURL == "" {
		log.Fatal("Missing Google OAuth configuration: ensure GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, and GOOGLE_REDIRECT_URL are set")
	}
}
