package src

import (
	"github.com/google/go-github/v53/github"
)

type CommentInfo struct {
	ID         int64
	OwnerName  string
	OwnerID    int64
	OwnerEmail string
	Repo       string
	PrNumber   int
	CreatedAt  github.Timestamp
	UpdatedAt  github.Timestamp
}

// PRInfo contains a subset of the github.PullRequest model payload
type PRInfo struct {
	OwnerName   string
	OwnerID     int64
	OwnerEmail  string
	Repo        string
	PrNumber    int
	CreatedAt   github.Timestamp
	UpdatedAt   github.Timestamp
	MergedAt    github.Timestamp
	CommentInfo []CommentInfo
}
