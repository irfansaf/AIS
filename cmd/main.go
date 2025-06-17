package main

import (
	"ais/internal/adapters/bitbucket"
	"ais/internal/adapters/http"
	"ais/internal/application"
	"log"
	"os"
	"time"

	"github.com/bytedance/sonic"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	webhookSecret := os.Getenv("BITBUCKET_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("WARNING: BITBUCKET_WEBHOOK_SECRET environment variable not set. Signature validation will be skipped.")
	}

	// Initialize Bitbucket client configuration
	bitbucketURL := os.Getenv("BITBUCKET_URL")
	bitbucketUsername := os.Getenv("BITBUCKET_USERNAME")
	bitbucketPassword := os.Getenv("BITBUCKET_PASSWORD")

	if bitbucketURL == "" {
		log.Fatal("BITBUCKET_URL environment variable is required")
	}
	if bitbucketUsername == "" {
		log.Fatal("BITBUCKET_USERNAME environment variable is required")
	}
	if bitbucketPassword == "" {
		log.Fatal("BITBUCKET_PASSWORD environment variable is required")
	}

	// Create Bitbucket client
	bitbucketClient := bitbucket.NewClient(bitbucketURL, bitbucketUsername, bitbucketPassword)

	// Create application service with dependencies
	appService := application.NewWebhookService(bitbucketClient)

	webhookHandler := http.NewWebhookHandler(appService, webhookSecret)

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"service":   "bitbucket-assistant",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/webhook/bitbucket", webhookHandler.HandlePullRequestEvent)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Bitbucket webhook listener on port %s...", port)
	log.Println("Endpoint available at POST http://localhost:8080/api/v1/webhook/bitbucket")

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start Fiber server: %v", err)
	}
}
