package main

import (
	"PUI/config"
	"PUI/db"
	"PUI/handlers"
	"PUI/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	//"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	config.LoadEnv()
	DB_MPK := db.Connect()
	DB_OPEN := db.ConnectOpen()

	r := gin.Default()

	/*
		// for use with token validator

		// Protect routes that require authentication
		authGroup := r.Group("/api")
		authGroup.Use(middleware.ValidateGoogleAccessToken())

		authGroup.GET("/user/info", handlers.GetUserInfo)
		authGroup.GET("/user/bets", handlers.GetUserBets)
	*/

	r.Use(middleware.DatabaseMiddlewareOpen(DB_OPEN))
	r.Use(middleware.DatabaseMiddlewareMPK(DB_MPK))

	// auth
	r.GET("/auth/google/login", handlers.HandleGoogleLogin)
	r.GET("/auth/google/callback", handlers.HandleGoogleCallback)

	// user
	r.GET("/api/user/me", handlers.GetUserInfo)

	// stops
	r.GET("/api/stops", handlers.GetStops)
	r.GET("/api/stop/:stopId", handlers.GetStopInfo)

	// trams
	r.GET("/api/tram/:tramId", handlers.GetTramInfo)

	// bets
	r.GET("/api/bets", handlers.GetUserBets)
	r.POST("/api/bets", handlers.PostBet)
	r.GET("/api/bets/:BetId", handlers.GetBetInfo)
	r.GET("/api/rates/:stopId", handlers.GetRates)

	err := r.Run(":8000")

	if err != nil {
		fmt.Println(err)
		return
	}
}
