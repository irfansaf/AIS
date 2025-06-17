package main

import (
	"ais/internal/adapters/http"
	"ais/internal/application"
	"github.com/bytedance/sonic"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	webhookSecret := os.Getenv("BITBUCKET_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("WARNING: BITBUCKET_WEBHOOK_SECRET environment variable not set. Signature validation will be skipped.")
	}

	appService := application.NewWebhookService()

	webhookHandler := http.NewWebhookHandler(appService, webhookSecret)

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
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
