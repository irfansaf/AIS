package bitbucket

import (
	"ais/internal/core/ports"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client implements the BitbucketClient interface
// This is a "driven adapter" in hexagonal architecture
type Client struct {
	baseURL    string
	username   string
	password   string // or token
	httpClient *http.Client
}

// NewClient creates a new Bitbucket client
func NewClient(baseURL, username, password string) ports.BitbucketClient {
	return &Client{
		baseURL:  baseURL,
		username: username,
		password: password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AddPullRequestComment implements the BitbucketClient interface
func (c *Client) AddPullRequestComment(ctx context.Context, projectKey, repoSlug string, pullRequestID int64, comment string) error {
	// Construct the API endpoint for adding comments
	// Bitbucket Server API: /rest/api/1.0/projects/{projectKey}/repos/{repositorySlug}/pull-requests/{pullRequestId}/comments
	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/pull-requests/%d/comments",
		c.baseURL, projectKey, repoSlug, pullRequestID)

	// Create the comment request payload
	commentReq := ports.CommentRequest{
		Text: comment,
	}

	// Marshal the request to JSON
	jsonData, err := json.Marshal(commentReq)
	if err != nil {
		return fmt.Errorf("failed to marshal comment request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set basic authentication
	req.SetBasicAuth(c.username, c.password)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bitbucket API returned status %d: %s", resp.StatusCode, resp.Status)
	}

	// Optionally, you can decode the response to get comment details
	var commentResp ports.CommentResponse
	if err := json.NewDecoder(resp.Body).Decode(&commentResp); err != nil {
		// Log the error but don't fail the operation since comment was created
		fmt.Printf("Warning: failed to decode comment response: %v\n", err)
	}

	return nil
}