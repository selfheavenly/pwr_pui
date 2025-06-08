package main

import (
	"PUI/config"
	"PUI/db"
	"PUI/handlers"
	"PUI/serialization"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/go-sql-driver/mysql"
	"log"
	//"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	config.LoadEnv()
	DB := db.Connect()

	r := gin.Default()

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

	jsonUsers, err := serialization.GetUsersJSON(DB)
	if err != nil {
		log.Fatal(err)
	}

	jsonBets, err := serialization.GetBetsJSON(DB)
	if err != nil {
		log.Fatal(err)
	}

	jsonStatusDictionary, err := serialization.GetStatusDictionaryJSON(DB)
	if err != nil {
		log.Fatal(err)
	}

	jsonStopStrainMap, err := serialization.GetStopTrainMapJSON(DB)
	if err != nil {
		log.Fatal(err)
	}

	jsonStopsDictionary, err := serialization.GetStopsDictionaryJSON(DB)
	if err != nil {
		log.Fatal(err)
	}

	jsonTramsDictionary, err := serialization.GetTramsDictionaryJSON(DB)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(jsonUsers)
	fmt.Println(jsonBets)
	fmt.Println(jsonStatusDictionary)
	fmt.Println(jsonStopStrainMap)
	fmt.Println(jsonStopsDictionary)
	fmt.Println(jsonTramsDictionary)

}
