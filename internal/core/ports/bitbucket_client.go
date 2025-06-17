package ports

import (
	"context"
)

// BitbucketClient defines the interface for interacting with Bitbucket API
// This is a "driven port" in hexagonal architecture - it defines what our application needs
// from external services without knowing how it's implemented
type BitbucketClient interface {
	// AddPullRequestComment adds a comment to a pull request
	AddPullRequestComment(ctx context.Context, projectKey, repoSlug string, pullRequestID int64, comment string) error
}

// CommentRequest represents the structure for creating a comment
type CommentRequest struct {
	Text string `json:"text"`
}

// CommentResponse represents the response from Bitbucket when creating a comment
type CommentResponse struct {
	ID          int64  `json:"id"`
	Version     int    `json:"version"`
	Text        string `json:"text"`
	Author      Author `json:"author"`
	CreatedDate int64  `json:"createdDate"`
	UpdatedDate int64  `json:"updatedDate"`
}

type Author struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}