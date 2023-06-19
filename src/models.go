package src

import (
	"github.com/google/go-github/v53/github"
)

type TeamType string

const (
	TeamBackend  TeamType = "backend"
	TeamFrontend TeamType = "frontend"
	TeamMobile   TeamType = "mobile"
	TeamDevOps   TeamType = "devops"
)

type CommentInfo struct {
	ID              int64
	OwnerName       string
	OwnerID         int64
	OwnerEmail      string
	Repo            string
	PrNumber        int
	contributorInfo SquadMember
	CreatedAt       github.Timestamp
	UpdatedAt       github.Timestamp
}

// PRInfo contains a subset of the github.PullRequest model payload
type PRInfo struct {
	OwnerName         string
	OwnerID           int64
	OwnerEmail        string
	Repo              string
	Team              TeamType
	PrNumber          int
	PrLink            string
	LinesAdded        int
	LinesDeleted      int
	TotalLinesChanged int
	NumOfComments     int
	contributorInfo   SquadMember
	CommentInfo       []CommentInfo
	CreatedAt         github.Timestamp
	UpdatedAt         github.Timestamp
	MergedAt          github.Timestamp
}

// SquadMember individual contributor info
type SquadMember struct {
	LoginName string
	SquadName string
	Team      TeamType
}
