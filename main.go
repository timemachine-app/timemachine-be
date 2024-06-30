package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/timemachine-app/timemachine-be/internal/config"
	"github.com/timemachine-app/timemachine-be/internal/handlers"
	"github.com/timemachine-app/timemachine-be/superbase"
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

	// Intialize Superbase
	superbaseClient := superbase.NewSupabaseClient(config.Clients.Superbase)

	// Initialize Router
	router := gin.Default()
	// Apply the rate limiting middleware
	router.Use(util.ValidationMiddleware(config.RateLimit, config.JwtSecret, superbaseClient))
	// health handler
	healthHandler := handlers.NewHealthHandler()
	router.GET("/health", healthHandler.IsHealthy)
	// account handler
	accountHandler := handlers.NewAccountHandler(config.Clients.SignInWithApple, superbaseClient, config.JwtSecret)
	router.POST("/signin/apple", accountHandler.SignInWithApple)
	router.POST("/delete", accountHandler.DeleteAccount)

	// event handler
	eventHandler := handlers.NewEventHandler(config.Clients.OpenAI, config.Clients.Gemini, config.Prompts.EventPrompts)
	router.POST("/event", eventHandler.ProcessEvent)
	router.POST("/search", eventHandler.Search)

	// setPortAndRun starts router on a server port
	router.Run(fmt.Sprintf(":%d", config.Server.Port))
}
