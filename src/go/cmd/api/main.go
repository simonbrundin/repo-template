package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Basic health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Readiness check (includes dependencies)
	r.GET("/ready", func(c *gin.Context) {
		checks := map[string]string{
			"api": "ok",
			"db":  "ok", // Lägg till faktisk DB-koll här när du har databas
		}

		status := http.StatusOK
		for name, check := range checks {
			if check != "ok" {
				status = http.StatusServiceUnavailable
				checks[name] = "error: " + check
			}
		}

		c.JSON(status, gin.H{
			"status":  "ready",
			"checks":  checks,
			"version": os.Getenv("APP_VERSION"),
		})
	})

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	r.Run(":8080")
}
