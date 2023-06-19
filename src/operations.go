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

	var lstOfMatchingPRs []PRInfo
	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return lstOfMatchingPRs, err
		}
		for _, pr := range prs {
			// skip draft
			if pr.GetDraft() {
				continue
			}
			// skip prs before from date and after that date
			createdAt := pr.GetCreatedAt()
			if createdAt.Before(from) || createdAt.After(to) {
				continue
			}
			//skip prs merged after to date
			mergedAt := pr.GetMergedAt()
			if mergedAt.After(to) {
				continue
			}

			contributor := contributors[pr.User.GetLogin()]
			prInfo := PRInfo{
				OwnerName:       pr.User.GetLogin(),
				OwnerEmail:      pr.User.GetEmail(),
				OwnerID:         pr.User.GetID(),
				Repo:            repo,
				PrNumber:        pr.GetNumber(),
				PrLink:          pr.GetHTMLURL(),
				contributorInfo: contributor,
				CreatedAt:       createdAt,
				UpdatedAt:       pr.GetUpdatedAt(),
				MergedAt:        mergedAt,
			}
			lstOfMatchingPRs = append(lstOfMatchingPRs, prInfo)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return lstOfMatchingPRs, nil
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

// GetPRLoc get pr lines of code additions and deletions
func GetPRLoc(
	ctx context.Context,
	org, repo string,
	prNumber int,
	client *github.Client,
) (additions, deletions int, err error) {

	opts := github.ListOptions{
		PerPage: 100,
	}
	for {
		files, resp, err := client.PullRequests.ListFiles(ctx, org, repo, prNumber, &opts)
		if err != nil {
			return 0, 0, err
		}
		for _, file := range files {
			additions += file.GetAdditions()
			deletions += file.GetDeletions()
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return
}
