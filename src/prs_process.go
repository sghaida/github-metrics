package src

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type RepoPrs struct {
	Repo      string
	OwnerTeam TeamType
	Prs       []PRInfo
}

type PRProcessor struct {
	wg           *sync.WaitGroup
	ch           chan<- RepoPrs
	client       *GithubClientsPool
	config       *Config
	contributors map[string]SquadMember
}

// NewPRProcessor create new PR processor
func NewPRProcessor(client *GithubClientsPool, conf *Config, ics map[string]SquadMember, send chan<- RepoPrs) *PRProcessor {
	return &PRProcessor{
		wg:           new(sync.WaitGroup),
		ch:           send,
		contributors: ics,
		client:       client,
		config:       conf,
	}
}

func (p *PRProcessor) process(repo string, team TeamType, counter *int32, from, to time.Time) {
	var enrichedPrs []PRInfo
	ctx := context.Background()
	client := p.client.Get()

	prs, err := GetRepoPrs(
		ctx, from, to, p.config.Org, repo, p.contributors, client,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, pr := range prs {
		comments, err := GetPrComments(ctx, p.config.Org, repo, pr.OwnerName, p.contributors, pr.PrNumber, client)
		if err != nil {
			fmt.Println(err)
			continue
		}
		additions, deletions, err := GetPRLoc(ctx, p.config.Org, repo, pr.PrNumber, client)
		if err != nil {
			fmt.Println(err)
			continue
		}

		pr.Team = team
		pr.CommentInfo = comments
		pr.LinesAdded = additions
		pr.LinesDeleted = deletions
		pr.TotalLinesChanged = additions + deletions
		pr.NumOfComments = len(pr.CommentInfo)
		enrichedPrs = append(enrichedPrs, pr)
	}

	p.ch <- RepoPrs{
		Repo:      repo,
		OwnerTeam: team,
		Prs:       enrichedPrs,
	}
	atomic.AddInt32(counter, 1)
}

func (p *PRProcessor) GetPrs(from, to time.Time) {
	var counter int32
	var reposCounter int
	atomic.StoreInt32(&counter, 0)

	for _, repos := range p.config.Repos {
		for _, repo := range repos.Names {
			go p.process(repo, repos.Type, &counter, from, to)
			reposCounter++
		}
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			c := atomic.LoadInt32(&counter)
			if c == int32(reposCounter) {
				close(p.ch)
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
