package main

import (
	"github.com/gin-gonic/gin"
	"gocheck/internal/handlers"
	"gocheck/internal/database"
	"gocheck/internal/middleware"
)

func main() {
	database.Connect()

	router := gin.Default()

	router.POST("/api/register", handlers.Register)
	router.POST("/api/login", handlers.Login)

	protected := router.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/sites", handlers.GetSites)
		protected.POST("/sites", handlers.CreateSite)
		protected.DELETE("/sites/:id", handlers.DeleteSite)
	}

	router.Run("localhost:8080")
}