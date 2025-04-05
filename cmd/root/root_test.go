package root

import (
	"os"
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	if cmd == nil {
		t.Fatal("NewRootCmd returned nil")
	}

	// Test command name
	if cmd.Use != "gus" {
		t.Errorf("Expected command name 'gus', got '%s'", cmd.Use)
	}

	// Test flags
	jsonFlag := cmd.Flags().Lookup("json")
	if jsonFlag == nil {
		t.Error("Expected 'json' flag to be defined")
	}
	if jsonFlag.Value.String() != "false" {
		t.Error("Expected 'json' flag to default to false")
	}

	pathFlag := cmd.Flags().Lookup("path")
	if pathFlag == nil {
		t.Error("Expected 'path' flag to be defined")
	}
	if pathFlag.Value.String() != "." {
		t.Error("Expected 'path' flag to default to '.'")
	}

	verboseFlag := cmd.Flags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("Expected 'verbose' flag to be defined")
	}
	if verboseFlag.Value.String() != "false" {
		t.Error("Expected 'verbose' flag to default to false")
	}
}

func TestRun(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "root-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test case 1: Default options
	cmd := NewRootCmd()
	rootPath = tempDir
	err = run(cmd, nil)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}

	// Test case 2: JSON output
	cmd = NewRootCmd()
	rootPath = tempDir
	jsonOutput = true
	err = run(cmd, nil)
	if err != nil {
		t.Errorf("Run failed with JSON output: %v", err)
	}

	// Test case 3: Verbose output
	cmd = NewRootCmd()
	rootPath = tempDir
	verbose = true
	err = run(cmd, nil)
	if err != nil {
		t.Errorf("Run failed with verbose output: %v", err)
	}

	// Test case 4: Invalid path
	cmd = NewRootCmd()
	rootPath = "/invalid/path"
	err = run(cmd, nil)
	if err == nil {
		t.Error("Expected error for invalid path")
	}

	// Test case 5: Path as argument
	cmd = NewRootCmd()
	rootPath = "."
	err = run(cmd, []string{tempDir})
	if err != nil {
		t.Errorf("Run failed with path argument: %v", err)
	}
}
