package http

import (
	"ais/internal/application"
	"ais/internal/core/domain/webhook"
	"ais/pkg/response"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"strings"
)

// WebhookHandler is the HTTP adapter. It holds dependencies like the application service.
type WebhookHandler struct {
	service     application.WebhookService
	secretToken string
}

// NewWebhookHandler creates a new WebhookHandler with the given service and secret token.
func NewWebhookHandler(service application.WebhookService, secretToken string) *WebhookHandler {
	return &WebhookHandler{
		service:     service,
		secretToken: secretToken,
	}
}

func (h *WebhookHandler) HandlePullRequestEvent(c *fiber.Ctx) error {
	rawBody := c.Body()

	if h.secretToken != "" {
		signatureHeader := c.Get("X-Hub-Signature")
		if err := h.validateSignature(signatureHeader, rawBody); err != nil {
			log.Errorf("Signature validation failed: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(response.NewErrorResponse("Invalid signature"))
		}
		log.Debugf("X-Hub-Signature validated successfully.")
	} else {
		log.Debugf("Warning: No secret token configured. Skipping signature validation.")
	}

	eventKey := c.Get("X-Event-Key")
	if eventKey == "" {
		log.Debugf("X-Event-Key header not found")
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse("Missing X-Event-Key header"))
	}

	if strings.HasPrefix(eventKey, "pr:") {
		var prEvent webhook.PullRequestEvent
		if err := sonic.Unmarshal(rawBody, &prEvent); err != nil {
			log.Errorf("Error unmarshalling PR event with Sonic: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse("Invalid JSON payload"))
		}

		if err := h.service.ProcessPullRequestEvent(c.Context(), prEvent); err != nil {
			log.Errorf("Error processing event '%s': %v", eventKey, err)
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse("Failed to process event"))
		}
	} else {
		log.Errorf("Received unhandled event type, acknowledging but not processing: %s", eventKey)
	}

	return c.SendStatus(fiber.StatusOK)
}

// validateSignature computes the HMAC of the raw body and compares it to the provided signature.
func (h *WebhookHandler) validateSignature(signatureHeader string, body []byte) error {
	if signatureHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing X-Hub-Signature header")
	}

	// The signature comes in the format "sha256=<hex_digest>"
	parts := strings.SplitN(signatureHeader, "=", 2)
	if len(parts) != 2 || parts[0] != "sha256" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid X-Hub-Signature format")
	}
	expectedSignature, err := hex.DecodeString(parts[1])
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid signature encoding")
	}

	mac := hmac.New(sha256.New, []byte(h.secretToken))
	mac.Write(body)
	calculatedSignature := mac.Sum(nil)

	if !hmac.Equal(calculatedSignature, expectedSignature) {
		return fiber.NewError(fiber.StatusUnauthorized, "Signature mismatch")
	}

	return nil
}
