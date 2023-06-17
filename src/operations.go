package src

import (
	"context"
	"github.com/google/go-github/v53/github"
	"time"
)

// GetRepoPrs get Pull requests based on from to date
func GetRepoPrs(
	ctx context.Context,
	from, to time.Time,
	owner, repo string,
	contributors map[string]SquadMember,
	client *github.Client,
) ([]PRInfo, error) {

	opts := &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	var allPRs []PRInfo
	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return allPRs, err
		}
		for _, pr := range prs {
			contributor := contributors[pr.User.GetLogin()]
			prInfo := PRInfo{
				OwnerName:       pr.User.GetLogin(),
				OwnerEmail:      pr.User.GetEmail(),
				OwnerID:         pr.User.GetID(),
				Repo:            repo,
				PrNumber:        *pr.Number,
				contributorInfo: contributor,
				CreatedAt:       pr.GetCreatedAt(),
				UpdatedAt:       pr.GetUpdatedAt(),
				MergedAt:        pr.GetMergedAt(),
			}
			allPRs = append(allPRs, prInfo)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Filter pull requests based on the date range
	var filteredPRs []PRInfo
	for _, pr := range allPRs {
		createdDate := pr.CreatedAt
		mergedAt := pr.MergedAt

		if createdDate.After(from) && createdDate.Before(to) && mergedAt.Before(to) {
			filteredPRs = append(filteredPRs, pr)
		}
	}
	return filteredPRs, nil
}

// GetPrComments return pr comments
func GetPrComments(
	ctx context.Context,
	org, repo, owner string,
	contributors map[string]SquadMember,
	prNumber int,
	client *github.Client) ([]CommentInfo, error) {

	opts := &github.PullRequestListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allComments []CommentInfo
	for {
		comments, resp, err := client.PullRequests.ListComments(
			ctx, org, repo, prNumber, nil,
		)
		if err != nil {
			return allComments, err
		}

		for _, comment := range comments {
			// skip pr owner comments
			if comment.User.GetLogin() == owner {
				continue
			}
			contributor := contributors[comment.User.GetLogin()]
			prComment := CommentInfo{
				ID:              comment.GetID(),
				OwnerName:       comment.User.GetLogin(),
				OwnerEmail:      comment.User.GetEmail(),
				OwnerID:         comment.User.GetID(),
				Repo:            repo,
				PrNumber:        prNumber,
				contributorInfo: contributor,
				CreatedAt:       comment.GetCreatedAt(),
				UpdatedAt:       comment.GetUpdatedAt(),
			}
			allComments = append(allComments, prComment)

		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allComments, nil
}
