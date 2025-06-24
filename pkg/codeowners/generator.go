package codeowners

import (
	"fmt"
	"os"
	"strings"

	"github.com/shraddha-aangiras/codeowners-generator/internal/github"
)

// GenerateCodeowners generates CODEOWNERS content from top contributors.
func GenerateCodeowners(contributors []github.Contributor, maxOwners int) string {
	var builder strings.Builder

	builder.WriteString("# Auto-generated CODEOWNERS file\n")
	builder.WriteString("# Do not edit manually\n\n")

	if maxOwners > len(contributors) {
		maxOwners = len(contributors)
	}

	for i := 0; i < maxOwners; i++ {
		builder.WriteString(fmt.Sprintf("* @%s\n", contributors[i].Login))
	}

	return builder.String()
}

// IsIdenticalToExisting checks if generated CODEOWNERS content matches what's already on disk.
func IsIdenticalToExisting(generated string, path string) (bool, error) {
	existing, err := os.ReadFile(path)
	if err != nil {
		// If file doesn't exist, it's not identical (but not an error either)
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return strings.TrimSpace(string(existing)) == strings.TrimSpace(generated), nil
}
