package formatter

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/nguyendangminh/gus/pkg/git"
)

func TestFormatRepositories(t *testing.T) {
	// Create test repositories
	repos := []*git.Repository{
		{
			Path: "/path/to/repo1",
			Changes: []string{
				"modified: file1.txt",
				"added: file2.txt",
			},
		},
		{
			Path: "/path/to/repo2",
			Changes: []string{
				"deleted: old_file.txt",
			},
		},
	}

	// Test case 1: Empty repositories
	err := FormatRepositories([]*git.Repository{}, FormatOptions{})
	if err != nil {
		t.Errorf("FormatRepositories failed with empty repos: %v", err)
	}

	// Test case 2: Text format
	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	err = FormatRepositories(repos, FormatOptions{})
	if err != nil {
		t.Errorf("FormatRepositories failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check output contains expected content
	expectedStrings := []string{
		"Found 2 Git repositories",
		"/path/to/repo1",
		"modified: file1.txt",
		"added: file2.txt",
		"/path/to/repo2",
		"deleted: old_file.txt",
	}

	for _, s := range expectedStrings {
		if !bytes.Contains([]byte(output), []byte(s)) {
			t.Errorf("Expected output to contain %q", s)
		}
	}

	// Test case 3: JSON format
	oldStdout = os.Stdout
	r, w, err = os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	err = FormatRepositories(repos, FormatOptions{JSON: true})
	if err != nil {
		t.Errorf("FormatRepositories failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read output
	buf.Reset()
	buf.ReadFrom(r)
	output = buf.String()

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}

	// Check JSON structure
	if _, ok := result["repositories"]; !ok {
		t.Error("Expected JSON output to have 'repositories' field")
	}
	if _, ok := result["metadata"]; !ok {
		t.Error("Expected JSON output to have 'metadata' field")
	}

	// Check metadata fields
	metadata, ok := result["metadata"].(map[string]interface{})
	if !ok {
		t.Error("Expected metadata to be an object")
	}

	if _, ok := metadata["scan_time"]; !ok {
		t.Error("Expected metadata to have 'scan_time' field")
	}
	if _, ok := metadata["total_repositories"]; !ok {
		t.Error("Expected metadata to have 'total_repositories' field")
	}
	if _, ok := metadata["version"]; !ok {
		t.Error("Expected metadata to have 'version' field")
	}

	// Check repositories
	repositories, ok := result["repositories"].([]interface{})
	if !ok {
		t.Error("Expected repositories to be an array")
	}
	if len(repositories) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(repositories))
	}

	// Check first repository
	repo1, ok := repositories[0].(map[string]interface{})
	if !ok {
		t.Error("Expected repository to be an object")
	}

	if repo1["path"] != "/path/to/repo1" {
		t.Errorf("Expected path to be '/path/to/repo1', got '%v'", repo1["path"])
	}

	changes1, ok := repo1["changes"].([]interface{})
	if !ok {
		t.Error("Expected changes to be an array")
	}
	if len(changes1) != 2 {
		t.Errorf("Expected 2 changes, got %d", len(changes1))
	}
}
