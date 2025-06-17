package application

import (
	"ais/internal/core/domain/webhook"
	"ais/internal/core/ports"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
)

// WebhookService defines the interface for our application's use cases.
// This is the "port" that driving adapters (like HTTP handlers) will call.
type WebhookService interface {
	ProcessPullRequestEvent(ctx context.Context, event webhook.PullRequestEvent) error
}

// webhookService is the concrete implementation of the WebhookService interface.
type webhookService struct {
	bitbucketClient ports.BitbucketClient
}

// NewWebhookService is a factory function to create a new webhookService.
// This is where you would inject its dependencies.
func NewWebhookService(bitbucketClient ports.BitbucketClient) WebhookService {
	return &webhookService{
		bitbucketClient: bitbucketClient,
	}
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
		
		// Add a welcome comment to the newly opened PR
		if err := s.addWelcomeComment(ctx, event); err != nil {
			log.Errorf("Failed to add welcome comment to PR #%d: %v", event.PullRequest.ID, err)
			// Don't return error as this is not critical to the webhook processing
		}

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

// addWelcomeComment adds a welcome comment to a newly opened pull request
func (s *webhookService) addWelcomeComment(ctx context.Context, event webhook.PullRequestEvent) error {
	projectKey := event.PullRequest.ToRef.Repository.Project.Key
	repoSlug := event.PullRequest.ToRef.Repository.Slug
	pullRequestID := event.PullRequest.ID

	// Create a personalized welcome message
	commentText := fmt.Sprintf(
		"🎉 **Welcome %s!**\n\n"+
			"Thank you for opening this pull request. Here are some quick reminders:\n\n"+
			"✅ **Code Review Checklist:**\n"+
			"- [ ] Code follows project coding standards\n"+
			"- [ ] Tests are included and passing\n"+
			"- [ ] Documentation is updated if needed\n"+
			"- [ ] No sensitive information is exposed\n\n"+
			"📋 **PR Details:**\n"+
			"- **Title:** %s\n"+
			"- **Source Branch:** `%s`\n"+
			"- **Target Branch:** `%s`\n\n"+
			"A reviewer will be assigned shortly. Thank you for your contribution! 🚀",
		event.PullRequest.Author.User.DisplayName,
		event.PullRequest.Title,
		event.PullRequest.FromRef.DisplayID,
		event.PullRequest.ToRef.DisplayID,
	)

	// Use the Bitbucket client to add the comment
	err := s.bitbucketClient.AddPullRequestComment(ctx, projectKey, repoSlug, pullRequestID, commentText)
	if err != nil {
		return fmt.Errorf("failed to add comment via Bitbucket API: %w", err)
	}

	log.Infof("Successfully added welcome comment to PR #%d in %s/%s", pullRequestID, projectKey, repoSlug)
	return nil
}
