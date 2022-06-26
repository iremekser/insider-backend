package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"project/app/api/handlers"
	"project/services/cache"
	"project/services/database"
	"time"
)

func main() {
	database.Init()
	cache.Init()

	// for web server operations
	// gin-gonic used
	router := gin.Default()

	// cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/get-fixture", handlers.GetFixture)
	router.POST("/start-next-week", handlers.StartOneWeek)
	router.POST("/start-all-week", handlers.StartAllWeek)
	router.POST("/clear-data", handlers.ClearData)
	router.GET("/get-teams", handlers.GetTeams)
	router.GET("/get-scores", handlers.GetScores)
	router.GET("/get-current-week", handlers.GetCurrentWeek)

	router.Run("localhost:8080")
}
