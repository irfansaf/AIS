package application

import (
	"ais/internal/core/domain/webhook"
	"context"
	"github.com/gofiber/fiber/v2/log"
)

// WebhookService defines the interface for our application's use cases.
// This is the "port" that driving adapters (like HTTP handlers) will call.
type WebhookService interface {
	ProcessPullRequestEvent(ctx context.Context, event webhook.PullRequestEvent) error
}

// webhookService is the concrete implementation of the WebhookService interface.
type webhookService struct {
	// In a real application, you would inject dependencies like repositories here.
	// For example:
	// prRepo repositories.PullRequestRepository
	// logger *log.Logger
}

// NewWebhookService is a factory function to create a new webhookService.
// This is where you would inject its dependencies.
func NewWebhookService() WebhookService {
	return &webhookService{}
}

// ProcessPullRequestEvent contains the business logic for handling a PR event.
// It is completely decoupled from how the event was received (e.g., HTTP).
func (s *webhookService) ProcessPullRequestEvent(ctx context.Context, event webhook.PullRequestEvent) error {
	log.Debugf("Received event: %s", event.EventKey)

	// Here, you orchestrate the business logic based on the event type.
	switch event.EventKey {
	case "pr:opened":
		log.Debugf("Pull Request #%d opened: '%s'", event.PullRequest.ID, event.PullRequest.Title)
		log.Debugf("Source Branch: %s -> Target Branch: %s", event.PullRequest.FromRef.DisplayID, event.PullRequest.ToRef.DisplayID)
		// TODO: Implement business logic for a newly opened PR.
		// Examples:
		// - Notify a Slack channel.
		// - Trigger a CI/CD pipeline.
		// - Add a default comment to the PR.

	case "pr:from_ref_updated":
		log.Debugf("Pull Request #%d source branch updated: '%s'", event.PullRequest.ID, event.PullRequest.Title)
		log.Debugf("New HEAD of source branch '%s' is %s", event.PullRequest.FromRef.DisplayID, event.PullRequest.FromRef.LatestCommit)
		log.Debugf("Previous HEAD was %s", event.PreviousFromHash)
		// TODO: Implement business logic for a PR update.
		// Examples:
		// - Re-run static analysis checks.
		// - Notify reviewers of the update.

	default:
		log.Debugf("Received unhandled pull request event type: %s", event.EventKey)
	}

	// In a real application, you might save state to a database via a repository interface.
	// err := s.prRepo.Save(ctx, &event.PullRequest)
	// if err != nil {
	//     return fmt.Errorf("failed to save pull request state: %w", err)
	// }

	return nil
}
