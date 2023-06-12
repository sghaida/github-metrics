package src

import (
	"context"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type Client struct {
	tokenSource oauth2.TokenSource
}

func NewClient(accessKey string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessKey},
	)
	return &Client{ts}
}

func (c *Client) Create(ctx context.Context) *github.Client {
	tc := oauth2.NewClient(ctx, c.tokenSource)
	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(tc.Transport)
	if err != nil {
		panic(err)
	}
	return github.NewClient(rateLimiter)
}
