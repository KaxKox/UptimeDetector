package main

import (
	"github.com/gin-gonic/gin"
	"gocheck/internal/handlers"
)

func main() {
	router := gin.Default()

	router.GET("/api/sites", handlers.GetSites)
	router.POST("/api/sites", handlers.CreateSite)
	router.PUT("/api/sites/:id", handlers.UpdateSite)
	router.DELETE("/api/sites/:id", handlers.DeleteSite)

	router.Run("localhost:8080")
}