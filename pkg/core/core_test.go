package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestScanner_Run(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "core-test-*")
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

	// Create a file in one of the repositories
	testFile := filepath.Join(tempDir, "dir1", "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test case 1: Default options
	options := Options{
		Path: tempDir,
	}
	scanner := New(options)
	err = scanner.Run()
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}

	// Test case 2: JSON output
	options = Options{
		Path: tempDir,
		JSON: true,
	}
	scanner = New(options)
	err = scanner.Run()
	if err != nil {
		t.Errorf("Run failed with JSON output: %v", err)
	}

	// Test case 3: Verbose output
	options = Options{
		Path:    tempDir,
		Verbose: true,
	}
	scanner = New(options)
	err = scanner.Run()
	if err != nil {
		t.Errorf("Run failed with verbose output: %v", err)
	}

	// Test case 4: Invalid path
	options = Options{
		Path: "/invalid/path",
	}
	scanner = New(options)
	err = scanner.Run()
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}
