package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsGitRepo(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test case 1: Not a Git repository
	if IsGitRepo(tempDir) {
		t.Error("Expected false for non-Git repository")
	}

	// Test case 2: Initialize a Git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize Git repository: %v", err)
	}

	if !IsGitRepo(tempDir) {
		t.Error("Expected true for Git repository")
	}

	// Test case 3: Non-existent directory
	if IsGitRepo("/non/existent/path") {
		t.Error("Expected false for non-existent directory")
	}
}

func TestCheckStatus(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize a Git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize Git repository: %v", err)
	}

	// Test case 1: Clean repository
	repo, err := CheckStatus(tempDir)
	if err != nil {
		t.Errorf("CheckStatus failed: %v", err)
	}
	if len(repo.Changes) != 0 {
		t.Errorf("Expected no changes, got %d changes", len(repo.Changes))
	}

	// Test case 2: Create a file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	repo, err = CheckStatus(tempDir)
	if err != nil {
		t.Errorf("CheckStatus failed: %v", err)
	}
	if len(repo.Changes) == 0 {
		t.Error("Expected changes after creating file")
	}

	// Test case 3: Add a file
	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	repo, err = CheckStatus(tempDir)
	if err != nil {
		t.Errorf("CheckStatus failed: %v", err)
	}
	if len(repo.Changes) == 0 {
		t.Error("Expected changes after adding file")
	}

	// Test case 4: Commit a file
	cmd = exec.Command("git", "commit", "-m", "Add test file")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit file: %v", err)
	}

	repo, err = CheckStatus(tempDir)
	if err != nil {
		t.Errorf("CheckStatus failed: %v", err)
	}
	if len(repo.Changes) != 0 {
		t.Errorf("Expected no changes after commit, got %d changes", len(repo.Changes))
	}

	// Test case 5: Modify a file
	if err := os.WriteFile(testFile, []byte("modified"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	repo, err = CheckStatus(tempDir)
	if err != nil {
		t.Errorf("CheckStatus failed: %v", err)
	}
	if len(repo.Changes) == 0 {
		t.Error("Expected changes after modifying file")
	}

	// Test case 6: Delete a file
	if err := os.Remove(testFile); err != nil {
		t.Fatalf("Failed to delete test file: %v", err)
	}

	repo, err = CheckStatus(tempDir)
	if err != nil {
		t.Errorf("CheckStatus failed: %v", err)
	}
	if len(repo.Changes) == 0 {
		t.Error("Expected changes after deleting file")
	}

	// Test case 7: Non-existent directory
	_, err = CheckStatus("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestParseGitStatus(t *testing.T) {
	// Test case 1: Empty output
	changes := parseGitStatus("")
	if len(changes) != 0 {
		t.Errorf("Expected no changes for empty output, got %d changes", len(changes))
	}

	// Test case 2: Modified file
	output := "M  file.txt"
	changes = parseGitStatus(output)
	if len(changes) != 1 {
		t.Errorf("Expected 1 change, got %d", len(changes))
	}
	if !strings.Contains(changes[0], "modified: file.txt") {
		t.Errorf("Expected 'modified: file.txt', got '%s'", changes[0])
	}

	// Test case 3: Added file
	output = "A  new_file.txt"
	changes = parseGitStatus(output)
	if len(changes) != 1 {
		t.Errorf("Expected 1 change, got %d", len(changes))
	}
	if !strings.Contains(changes[0], "added: new_file.txt") {
		t.Errorf("Expected 'added: new_file.txt', got '%s'", changes[0])
	}

	// Test case 4: Deleted file
	output = "D  old_file.txt"
	changes = parseGitStatus(output)
	if len(changes) != 1 {
		t.Errorf("Expected 1 change, got %d", len(changes))
	}
	if !strings.Contains(changes[0], "deleted: old_file.txt") {
		t.Errorf("Expected 'deleted: old_file.txt', got '%s'", changes[0])
	}

	// Test case 5: Multiple changes
	output = "M  file1.txt\nA  file2.txt\nD  file3.txt"
	changes = parseGitStatus(output)
	if len(changes) != 3 {
		t.Errorf("Expected 3 changes, got %d", len(changes))
	}
}
