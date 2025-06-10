package main

import (
	"PUI/config"
	"PUI/db"
	"PUI/handlers"
	"PUI/middleware"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	config.LoadEnv()
	DB_MPK := db.Connect()
	DB_OPEN := db.ConnectOpen()

	r := gin.Default()

	r.Use(middleware.DatabaseMiddlewareOpen(DB_OPEN))
	r.Use(middleware.DatabaseMiddlewareMPK(DB_MPK))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// auth
	r.GET("/auth/google/login", handlers.HandleGoogleLogin)
	r.GET("/auth/google/callback", handlers.HandleGoogleCallback)

	// Apply middleware to routes that require authentication
	protected := r.Group("/api")
	protected.Use(middleware.GoogleIDMiddleware())
	{
		// user
		protected.GET("/user/me", handlers.GetUserInfo)

		// stops
		protected.GET("/stops", handlers.GetStops)
		protected.GET("/stop/:stopId", handlers.GetStopInfo)

		// trams
		protected.GET("/tram/:tramId", handlers.GetTramInfo)

		// bets
		protected.GET("/bets", handlers.GetUserBets)
		protected.POST("/bets", handlers.PostBet)
		protected.GET("/bets/:betId", handlers.GetBetInfo)
		protected.GET("/rates/:stopId", handlers.GetRates)
	}

	/*
		// for use with token validator

		// Protect routes that require authentication
		authGroup := r.Group("/api")
		authGroup.Use(middleware.ValidateGoogleAccessToken())

		authGroup.GET("/user/info", handlers.GetUserInfo)
		authGroup.GET("/user/bets", handlers.GetUserBets)
	*/

	err := r.Run(":8000")

	if err != nil {
		fmt.Println(err)
		return
	}
}
