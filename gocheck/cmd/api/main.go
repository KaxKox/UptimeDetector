package main

import (
	"gocheck/internal/database"
	"gocheck/internal/handlers"
	"gocheck/internal/middleware"
	"gocheck/internal/monitor"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	go monitor.Start()

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.Use(middleware.APILogger())

	router.POST("/api/register", handlers.Register)
	router.POST("/api/login", handlers.Login)

	router.GET("/ws", handlers.WsHandler)

	protected := router.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/sites", handlers.GetSites)
		protected.POST("/sites", handlers.CreateSite)
		protected.DELETE("/sites/:id", handlers.DeleteSite)
		protected.PUT("/sites/:id", handlers.UpdateSite)

		protected.GET("/export", handlers.ExportCSV)
		protected.POST("/import", handlers.ImportCSV)
	}

	router.Run("localhost:8080")
	//router.RunTLS("localhost:8080", "cert.pem", "key.pem")
}
