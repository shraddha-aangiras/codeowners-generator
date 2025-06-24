package main

import (
	"os"
	"fmt"
	"time"
	"github.com/shraddha-aangiras/codeowners-generator/internal/github"
	"github.com/shraddha-aangiras/codeowners-generator/pkg/utils"
	"github.com/shraddha-aangiras/codeowners-generator/pkg/codeowners"
)

func main() {
	// Parse CLI args
	githubServerURL, organizationName, repositoryName, githubToken, duration, codeReviewersCount, err := utils.ParseArgs()
	if err != nil {
		fmt.Println("Error parsing CLI arguments:", err)
		return
	}

	// Print parsed args for debugging
	fmt.Println("Parsed CLI args:")
	fmt.Println("  GitHub URL:   ", githubServerURL)
	fmt.Println("  Org:          ", organizationName)
	fmt.Println("  Repo:         ", repositoryName)
	fmt.Println("  Token present:", githubToken != "")

	// Initialize GitHub client
	client := github.NewClient(
		githubToken,
		githubServerURL,
		organizationName,
		repositoryName,
	)

	// Fetch contributors
	since := time.Now().Add(-duration) // where duration is your time.Duration variable
	contributors, err := client.GetTopContributors(since)

	if err != nil {
		fmt.Println("Error fetching contributors:", err)
		return
	}

	if len(contributors) == 0 {
		fmt.Println("No contributors found.")
		return
	}

	// Limit to top N contributors
	if len(contributors) > codeReviewersCount {
		contributors = contributors[:codeReviewersCount]
	}

	// Output contributors
	fmt.Println("\nTop contributors:")
	for _, c := range contributors {
		fmt.Printf("- @%s (%d commits)\n", c.Login, c.Contributions)
	}

	if len(contributors) == 0 {
		fmt.Println("No contributors found. Try increasing --duration or check if the repo is active.")
		return
	}
	
	if len(contributors) < codeReviewersCount {
		fmt.Printf("Only %d contributor(s) found, but you requested %d code owners.\n", len(contributors), codeReviewersCount)
		fmt.Println("Proceeding with available contributors.")
		codeReviewersCount = len(contributors)
	}
	
	fmt.Printf("Generating CODEOWNERS with %d top contributor(s)...\n", codeReviewersCount)
	codeownersContent := codeowners.GenerateCodeowners(contributors, codeReviewersCount)
	
	identical, err := codeowners.IsIdenticalToExisting(codeownersContent, "CODEOWNERS")
	if err != nil {
		fmt.Println("Failed to compare with existing CODEOWNERS file:", err)
		return
	}
	
	if identical {
		fmt.Println("CODEOWNERS file is already up to date — no changes needed.")
		return
	}
	
	err = os.WriteFile("CODEOWNERS", []byte(codeownersContent), 0644)
	if err != nil {
		fmt.Println("Failed to write CODEOWNERS file:", err)
		return
	}
	
	fmt.Println("CODEOWNERS file written successfully!")
	
}
