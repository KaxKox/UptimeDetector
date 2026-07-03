package main

import (
	"github.com/gin-gonic/gin"
	"gocheck/internal/handlers"
	"gocheck/internal/database"
)

func main() {
	database.Connect()

	router := gin.Default()

	router.GET("/api/sites", handlers.GetSites)
	router.POST("/api/sites", handlers.CreateSite)
	router.DELETE("/api/sites/:id", handlers.DeleteSite)

	router.Run("localhost:8080")
}