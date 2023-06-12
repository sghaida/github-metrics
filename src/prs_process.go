package src

import (
	"context"
	"fmt"
	"github.com/google/go-github/v53/github"
	"sync"
	"sync/atomic"
	"time"
)

type RepoPrs struct {
	Repo string
	Prs  []PRInfo
}

type PRProcessor struct {
	wg     *sync.WaitGroup
	ch     chan<- RepoPrs
	client *github.Client
	config *Config
}

// NewPRProcessor create new PR processor
func NewPRProcessor(client *github.Client, conf *Config, send chan<- RepoPrs) *PRProcessor {
	return &PRProcessor{
		wg:     new(sync.WaitGroup),
		ch:     send,
		client: client,
		config: conf,
	}
}

func (p *PRProcessor) process(owner, repo string, counter *int32, from, to time.Time) {
	var enrichedPrs []PRInfo
	ctx := context.Background()

	prs, err := GetRepoPrs(
		ctx, from, to, p.config.Org, repo, p.client,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, pr := range prs {
		comments, err := GetPrComments(ctx, p.config.Org, repo, pr.OwnerName, pr.PrNumber, p.client)
		if err != nil {
			fmt.Println(err)
			continue
		}
		pr.CommentInfo = comments
		enrichedPrs = append(enrichedPrs, pr)
	}

	p.ch <- RepoPrs{
		Repo: repo,
		Prs:  enrichedPrs,
	}
	atomic.AddInt32(counter, 1)
}

func (p *PRProcessor) GetPrs(from, to time.Time) {
	var counter int32
	atomic.StoreInt32(&counter, 0)

	for _, repo := range p.config.Repos.Backend {
		go p.process(config.Org, repo, &counter, from, to)
	}
	for _, repo := range p.config.Repos.Frontend {
		go p.process(config.Org, repo, &counter, from, to)
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			c := atomic.LoadInt32(&counter)
			if c == int32(len(config.Repos.Backend)+len(config.Repos.Frontend)) {
				close(p.ch)
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()

}
