package github

import (
	"context"
	"sort"
	"time"
	"fmt"

	"github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"
)

// Contributor holds minimal info about a contributor and their commit count
type Contributor struct {
	Login         string
	Contributions int
}

type Client struct {
	client *github.Client
	owner  string
	repo   string
}

// NewClient initializes a GitHub API client with auth token
func NewClient(token, githubServerURL, owner, repo string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &Client{
		client: client,
		owner:  owner,
		repo:   repo,
	}
}

// GetTopContributors gets the top contributors for given time duration
func (c *Client) GetTopContributors(since time.Time) ([]Contributor, error) {
	ctx := context.Background()

	commitCount := make(map[string]int)

	opt := &github.CommitsListOptions{
		Since: since,
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1, 
		},
	}

	for {
		commits, resp, err := c.client.Repositories.ListCommits(ctx, c.owner, c.repo, opt)
		if err != nil {
			return nil, err
		}


		for _, commit := range commits {
			commitDate := commit.GetCommit().GetAuthor().GetDate()

			// extra (solves possible issue number discrepancy?)
			if commitDate.Before(since) {
				fmt.Printf("Skipping commit %s dated %v (before %v)\n", commit.GetSHA(), commitDate, since)
				continue
			}

			if commit.Author != nil && commit.Author.Login != nil {
				login := *commit.Author.Login
				commitCount[login]++
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// Convert map to slice
	contributors := make([]Contributor, 0, len(commitCount))
	for login, count := range commitCount {
		contributors = append(contributors, Contributor{
			Login:        login,
			Contributions: count,
		})
	}

	// Sort descending by contributions
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Contributions > contributors[j].Contributions
	})

	return contributors, nil
}
