package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	openapiui "github.com/PeterTakahashi/gin-openapi/openapiui"
)

// @title			Repo Template API
// @version			1.0
// @description		API for repo-template project
// @host			localhost:8080
// @BasePath		/

func getDocsPath() string {
	if path := os.Getenv("DOCS_PATH"); path != "" {
		return path
	}
	return "./cmd/api/docs/swagger.json"
}

func main() {
	r := gin.Default()

	// API documentation with Scalar UI
	r.GET("/docs/*any", openapiui.WrapHandler(openapiui.Config{
		SpecURL:      "/docs/openapi.json",
		SpecFilePath: getDocsPath(),
		Title:        "Repo Template API",
		Theme:        "dark",
	}))

	// @Summary		Health check
	// @Description	Returns the health status of the API
	// @Tags			health
	// @Produce		json
	// @Success		200	{object}	map[string]interface{}
	// @Router		/health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"version":   os.Getenv("APP_VERSION"),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// @Summary		Hello world
	// @Description	Returns a hello message
	// @Tags			greeting
	// @Produce		json
	// @Success		200	{object}	map[string]string
	// @Router		/hello [get]
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	r.Run(":8080")
}
