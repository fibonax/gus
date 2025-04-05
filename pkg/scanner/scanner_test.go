package scanner

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestScan(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scanner-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	dirs := []string{
		"dir1",
		"dir1/subdir1",
		"dir2",
		"dir2/subdir2",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tempDir, dir), 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Initialize Git repositories in some directories
	gitDirs := []string{
		"dir1",
		"dir2/subdir2",
	}

	for _, dir := range gitDirs {
		cmd := exec.Command("git", "init")
		cmd.Dir = filepath.Join(tempDir, dir)
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to initialize Git repository in %s: %v", dir, err)
		}
	}

	// Test scanning
	s := New(tempDir)
	gitDirsFound, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Check results
	if len(gitDirsFound) != len(gitDirs) {
		t.Errorf("Expected %d Git repositories, found %d", len(gitDirs), len(gitDirsFound))
	}

	// Check if all expected Git repositories were found
	for _, expectedDir := range gitDirs {
		found := false
		expectedPath := filepath.Join(tempDir, expectedDir)
		for _, foundDir := range gitDirsFound {
			if foundDir == expectedPath {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find Git repository in %s", expectedDir)
		}
	}
}
