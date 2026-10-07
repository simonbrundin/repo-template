package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	openapiui "github.com/PeterTakahashi/gin-openapi/openapiui"
)

// @title			Repo Template API
// @version			1.0
// @description		API for repo-template project
// @BasePath		/

func getDocsPath() string {
	if path := os.Getenv("DOCS_PATH"); path != "" {
		return path
	}
	return "./cmd/api/docs/swagger.json"
}

// readSwaggerSpec reads and returns the swagger spec
func readSwaggerSpec() ([]byte, error) {
	return os.ReadFile(getDocsPath())
}

// withHost modifies the swagger spec to use the correct host
func withHost(specJSON []byte, host string) ([]byte, error) {
	var spec map[string]interface{}
	if err := json.Unmarshal(specJSON, &spec); err != nil {
		return nil, err
	}
	spec["host"] = host
	return json.Marshal(spec)
}

func main() {
	r := gin.Default()

	// API documentation with Scalar UI
	r.GET("/docs/*any", openapiui.WrapHandler(openapiui.Config{
		SpecURL: "/docs/openapi.json",
		SpecProvider: func() ([]byte, error) {
			host := os.Getenv("API_HOST")
			if host == "" {
				host = "localhost:8080" // fallback
			}
			spec, err := readSwaggerSpec()
			if err != nil {
				return nil, err
			}
			return withHost(spec, host)
		},
		Title: "Repo Template API",
		Theme: "dark",
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

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" || !strings.HasPrefix(port, ":") {
		if port == "" {
			port = "8080"
		}
		port = ":" + port
	}

	r.Run(port)
}
