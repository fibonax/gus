package scanner

import (
	"os"
	"path/filepath"

	"github.com/nguyendangminh/gus/pkg/git"
)

// Scanner represents a directory scanner
type Scanner struct {
	rootPath string
}

// New creates a new Scanner instance
func New(rootPath string) *Scanner {
	return &Scanner{
		rootPath: rootPath,
	}
}

// Scan performs a recursive scan of the root directory
func (s *Scanner) Scan() ([]string, error) {
	var gitDirs []string

	err := filepath.Walk(s.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a directory
		if !info.IsDir() {
			return nil
		}

		// Check if this is a Git repository
		if git.IsGitRepo(path) {
			gitDirs = append(gitDirs, path)
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return gitDirs, nil
}
