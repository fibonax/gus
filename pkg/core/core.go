package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nguyendangminh/gus/pkg/formatter"
	"github.com/nguyendangminh/gus/pkg/git"
	"github.com/nguyendangminh/gus/pkg/scanner"
)

// Options contains all options for scanning
type Options struct {
	Path    string
	JSON    bool
	Verbose bool
}

// Scanner represents the main scanner
type Scanner struct {
	options Options
}

// New creates a new Scanner instance
func New(options Options) *Scanner {
	return &Scanner{
		options: options,
	}
}

// Run performs the scanning process
func (s *Scanner) Run() error {
	// Convert to absolute path
	absPath, err := filepath.Abs(s.options.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", absPath)
	}

	// Create directory scanner
	dirScanner := scanner.New(absPath)
	gitDirs, err := dirScanner.Scan()
	if err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	if s.options.Verbose {
		fmt.Printf("Found %d Git repositories\n", len(gitDirs))
	}

	// Check status of each repository
	var reposWithChanges []*git.Repository
	for _, dir := range gitDirs {
		repo, err := git.CheckStatus(dir)
		if err != nil {
			if s.options.Verbose {
				fmt.Fprintf(os.Stderr, "Warning: failed to check status of %s: %v\n", dir, err)
			}
			continue
		}

		if len(repo.Changes) > 0 {
			reposWithChanges = append(reposWithChanges, repo)
		}
	}

	// Format and print results
	opts := formatter.FormatOptions{
		JSON: s.options.JSON,
	}
	return formatter.FormatRepositories(reposWithChanges, opts)
}
