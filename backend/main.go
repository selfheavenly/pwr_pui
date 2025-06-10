package main

import (
	"PUI/config"
	"PUI/db"
	"PUI/handlers"
	"PUI/middleware"
	"fmt"
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

	// auth
	r.GET("/auth/google/login", handlers.HandleGoogleLogin)
	r.GET("/auth/google/callback", handlers.HandleGoogleCallback)

	// Apply middleware to routes that require authentication

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/stops", handlers.GetStops)
		apiGroup.GET("/stop/:stopId", handlers.GetStopInfo)

		// bets
		apiGroup.GET("/rates/:stopId", handlers.GetRates)
	}

	protected := r.Group("/api")
	protected.Use(middleware.ValidateGoogleAccessToken())
	{
		protected.GET("/user/me", handlers.GetUserInfo) // git
		protected.GET("/bets", handlers.GetUserBets)
		protected.POST("/bets", handlers.PostBet) // git
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
