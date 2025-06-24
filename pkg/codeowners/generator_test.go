package codeowners

import (
	"os"
	"testing"

	"github.com/shraddha-aangiras/codeowners-generator/internal/github"
)

func TestGenerateCodeowners(t *testing.T) {
	contributors := []github.Contributor{
		{Login: "alice"},
		{Login: "bob"},
		{Login: "carol"},
	}

	// Test normal case: maxOwners less than contributors count
	content := GenerateCodeowners(contributors, 2)
	expected := "# Auto-generated CODEOWNERS file\n# Do not edit manually\n\n* @alice\n* @bob\n"
	if content != expected {
		t.Errorf("expected:\n%s\nbut got:\n%s", expected, content)
	}

	// Test maxOwners more than contributors count (should not panic, just list all)
	content = GenerateCodeowners(contributors, 10)
	expected = "# Auto-generated CODEOWNERS file\n# Do not edit manually\n\n* @alice\n* @bob\n* @carol\n"
	if content != expected {
		t.Errorf("expected:\n%s\nbut got:\n%s", expected, content)
	}
}

func TestIsIdenticalToExisting(t *testing.T) {
	tmpFile := "test_CODEOWNERS"
	defer os.Remove(tmpFile) // clean up

	content := "# Auto-generated CODEOWNERS file\n* @alice\n"

	// Write content to file
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Should return true when comparing identical content
	identical, err := IsIdenticalToExisting(content, tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !identical {
		t.Errorf("expected identical to be true")
	}

	// Should return false for different content
	identical, err = IsIdenticalToExisting(content+"more", tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identical {
		t.Errorf("expected identical to be false")
	}

	// Test file not existing (should return false, no error)
	identical, err = IsIdenticalToExisting(content, "nonexistent_file")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identical {
		t.Errorf("expected identical to be false when file missing")
	}
}
