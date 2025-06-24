package utils

import (
	"fmt"
	"os"
	"time"
	"github.com/urfave/cli/v2"
)

func ParseArgs() (string, string, string, string, time.Duration, int, error) {
	var (
		githubServerURL    string
		organizationName   string
		repositoryName     string
		githubToken        string
		duration           time.Duration
		codeReviewersCount int
	)

	app := &cli.App{
		Name:  "codeowners-generator",
		Usage: "Generate CODEOWNERS file based on GitHub contributor activity",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "github-server-url",
				EnvVars: []string{"GITHUB_SERVER_URL"},
				Value:   "https://api.github.com",
				Usage:   "GitHub API server URL",
			},
			&cli.StringFlag{
				Name:     "organization-name",
				EnvVars:  []string{"ORGANIZATION_NAME"},
				Usage:    "GitHub organization name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "repository-name",
				EnvVars:  []string{"REPOSITORY_NAME"},
				Usage:    "GitHub repository name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "github-token",
				EnvVars:  []string{"GITHUB_TOKEN"},
				Usage:    "GitHub personal access token",
				Required: true,
			},
			&cli.DurationFlag{
				Name:    "duration",
				EnvVars: []string{"DURATION"},
				Value:   time.Hour * 24 * 30,
				Usage:   "Duration to look back for contributions (e.g., 720h = 30d)",
			},
			&cli.IntFlag{
				Name:    "code-reviewers-count",
				EnvVars: []string{"CODE_REVIEWERS_COUNT"},
				Value:   3,
				Usage:   "Number of top contributors to assign as code owners",
			},
		},

		Action: func(c *cli.Context) error {
			githubServerURL = c.String("github-server-url")
			organizationName = c.String("organization-name")
			repositoryName = c.String("repository-name")
			githubToken = c.String("github-token")
			duration = c.Duration("duration")
			codeReviewersCount = c.Int("code-reviewers-count")
			return nil
		},
	}
	
	err := app.Run(os.Args)
	if err != nil {
		return "", "", "", "", 0, 0, fmt.Errorf("CLI parse failed: %w", err)
	}

	return githubServerURL, organizationName, repositoryName, githubToken, duration, codeReviewersCount, nil
}


