package src

import (
	"context"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
	"sync/atomic"
)

type GithubClientsPool struct {
	clients  []*github.Client
	balancer int32
}

func NewGithubClientsPool(accessKeys []string) *GithubClientsPool {
	var clients []*github.Client
	for _, accessKey := range accessKeys {
		tokenSource := newTokenSource(accessKey)
		client := tokenSource.createGithubClient(context.Background())
		clients = append(clients, client)
	}
	return &GithubClientsPool{clients: clients}
}

// Get return one of the clients based on round-robin balancing approach
func (cp *GithubClientsPool) Get() *github.Client {
	value := atomic.AddInt32(&cp.balancer, 1)
	index := int(value) % len(cp.clients)
	return cp.clients[index]
}

type TokenSource struct {
	tokenSource oauth2.TokenSource
}

func newTokenSource(accessKey string) *TokenSource {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessKey},
	)
	return &TokenSource{tokenSource: ts}
}

func (c *TokenSource) createGithubClient(ctx context.Context) *github.Client {
	tc := oauth2.NewClient(ctx, c.tokenSource)
	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(tc.Transport)
	if err != nil {
		panic(err)
	}
	return github.NewClient(rateLimiter)
}
