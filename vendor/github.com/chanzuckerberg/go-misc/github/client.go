package github

import (
	"net/http"

	"github.com/google/go-github/github"
)

// Client is a github client
type Client struct {
	client *github.Client
}

// NewClient returns a new github client
func NewClient(httpClient *http.Client) *Client {
	githubClient := github.NewClient(httpClient)

	return &Client{
		client: githubClient,
	}
}
