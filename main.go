package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/timemachine-app/timemachine-be/internal/config"
	"github.com/timemachine-app/timemachine-be/internal/handlers"
	"github.com/timemachine-app/timemachine-be/util"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "default"
	}

	// Initialize configuration
	config, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Router
	router := gin.Default()
	// Apply the rate limiting middleware
	router.Use(util.RateLimitMiddleware(config.RateLimit))
	// health handler
	healthHandler := handlers.NewHealthHandler()
	router.GET("/health", healthHandler.IsHealthy)
	// event handler
	eventHandler := handlers.NewEventHandler(config.Clients.OpenAI, config.Prompts.EventPrompts)
	router.POST("/event", eventHandler.ProcessEvent)
	router.POST("/search", eventHandler.Search)

	// setPortAndRun starts router on a server port
	router.Run(fmt.Sprintf(":%d", config.Server.Port))
}
