package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

// parseFlexibleDuration parses durations like "1w2d3h" into time.Duration
func parseFlexibleDuration(input string) (time.Duration, error) {
	re := regexp.MustCompile(`(?i)(\d+)([wdhms])`)
	matches := re.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration format: %s", input)
	}

	var total time.Duration
	for _, match := range matches {
		valStr, unit := match[1], strings.ToLower(match[2])
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return 0, err
		}

		switch unit {
		case "w":
			total += time.Duration(val) * 7 * 24 * time.Hour
		case "d":
			total += time.Duration(val) * 24 * time.Hour
		case "h":
			total += time.Duration(val) * time.Hour
		case "m":
			total += time.Duration(val) * time.Minute
		case "s":
			total += time.Duration(val) * time.Second
		default:
			return 0, fmt.Errorf("unknown time unit: %s", unit)
		}
	}

	return total, nil
}

func ParseArgs() (string, string, string, string, time.Duration, int, error) {
	var (
		githubServerURL    string
		organizationName   string
		repositoryName     string
		githubToken        string
		durationInput      string
		parsedDuration     time.Duration
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
			&cli.StringFlag{
				Name:    "duration",
				EnvVars: []string{"DURATION"},
				Value:   "30d", // default = 30 days
				Usage:   "Lookback duration for commits (e.g. 2w3d5h)",
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
			durationInput = c.String("duration")
			codeReviewersCount = c.Int("code-reviewers-count")

			// Custom parsing
			var err error
			parsedDuration, err = parseFlexibleDuration(durationInput)
			if err != nil {
				return cli.Exit(fmt.Sprintf("invalid value %q for flag -duration: parse error", durationInput), 1)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return "", "", "", "", 0, 0, fmt.Errorf("CLI parse failed: %w", err)
	}

	return githubServerURL, organizationName, repositoryName, githubToken, parsedDuration, codeReviewersCount, nil
}