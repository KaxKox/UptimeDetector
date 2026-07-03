package middleware

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
)

func APILogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		log.Printf("STATYSTYKI API %s | %s %s | Status: %d | Czas: %v", clientIP, method, path, statusCode, duration)
	}
}