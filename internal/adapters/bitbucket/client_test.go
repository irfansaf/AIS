package bitbucket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_AddPullRequestComment(t *testing.T) {
	// Create a test server to mock Bitbucket API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request method and path
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		
		expectedPath := "/rest/api/1.0/projects/TEST/repos/test-repo/pull-requests/123/comments"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		
		// Verify headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		
		// Verify basic auth
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Expected basic auth to be present")
		}
		if username != "testuser" || password != "testpass" {
			t.Errorf("Expected auth testuser:testpass, got %s:%s", username, password)
		}
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": 1,
			"version": 1,
			"text": "Test comment",
			"author": {
				"name": "testuser",
				"displayName": "Test User",
				"emailAddress": "test@example.com"
			},
			"createdDate": 1234567890000,
			"updatedDate": 1234567890000
		}`))
	}))
	defer testServer.Close()
	
	// Create client with test server URL
	client := NewClient(testServer.URL, "testuser", "testpass")
	
	// Test the AddPullRequestComment method
	err := client.AddPullRequestComment(
		context.Background(),
		"TEST",
		"test-repo",
		123,
		"Test comment",
	)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClient_AddPullRequestComment_Error(t *testing.T) {
	// Create a test server that returns an error
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errors":[{"message":"Authentication failed"}]}`))
	}))
	defer testServer.Close()
	
	// Create client with test server URL
	client := NewClient(testServer.URL, "testuser", "wrongpass")
	
	// Test the AddPullRequestComment method with error
	err := client.AddPullRequestComment(
		context.Background(),
		"TEST",
		"test-repo",
		123,
		"Test comment",
	)
	
	if err == nil {
		t.Error("Expected error, got nil")
	}
}