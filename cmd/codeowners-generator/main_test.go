package main

import (
	"os"
	"testing"
	"time"

	"github.com/shraddha-aangiras/codeowners-generator/internal/github"
	"github.com/shraddha-aangiras/codeowners-generator/pkg/codeowners"
)

// MockGitHubClient implements the same method as your real client
type MockGitHubClient struct {
	Contributors []github.Contributor
	Err          error
}

func (m *MockGitHubClient) GetTopContributors(since time.Time) ([]github.Contributor, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Contributors, nil
}

func TestGenerateCodeownersFromMock(t *testing.T) {
	mockClient := &MockGitHubClient{
		Contributors: []github.Contributor{
			{Login: "alice", Contributions: 10},
			{Login: "bob", Contributions: 5},
		},
	}

	expectedContent := `# Auto-generated CODEOWNERS file
# Do not edit manually

* @alice
* @bob
`

	content := codeowners.GenerateCodeowners(mockClient.Contributors, 2)

	// Clean up before and after
	_ = os.Remove("CODEOWNERS")
	defer os.Remove("CODEOWNERS")

	// Write CODEOWNERS file
	err := os.WriteFile("CODEOWNERS", []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write CODEOWNERS file: %v", err)
	}

	// Read it back
	written, err := os.ReadFile("CODEOWNERS")
	if err != nil {
		t.Fatalf("Failed to read CODEOWNERS file: %v", err)
	}

	if string(written) != expectedContent {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedContent, string(written))
	}
}

func TestTooFewContributorsFallback(t *testing.T) {
	mockClient := &MockGitHubClient{
		Contributors: []github.Contributor{
			{Login: "solo-dev", Contributions: 1},
		},
	}

	// Ask for more code owners than available
	expectedContent := `# Auto-generated CODEOWNERS file
# Do not edit manually

* @solo-dev
`

	content := codeowners.GenerateCodeowners(mockClient.Contributors, 3)

	if content != expectedContent {
		t.Errorf("Expected fallback content:\n%s\nGot:\n%s", expectedContent, content)
	}
}
