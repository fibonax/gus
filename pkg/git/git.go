package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Repository represents a Git repository
type Repository struct {
	Path    string
	Changes []string
}

// IsGitRepo checks if a directory is a Git repository
func IsGitRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// NewRepository creates a new Repository instance
func NewRepository(path string) *Repository {
	return &Repository{
		Path: path,
	}
}

// CheckStatus checks the status of a Git repository
func CheckStatus(repoPath string) (*Repository, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		Path:    repoPath,
		Changes: parseGitStatus(string(output)),
	}

	return repo, nil
}

// parseGitStatus parses the output of git status --porcelain
func parseGitStatus(output string) []string {
	if output == "" {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	changes := make([]string, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse the status line
		// Format: XY PATH
		// X = staged changes
		// Y = unstaged changes
		// Both can be:
		//   M = modified
		//   A = added
		//   D = deleted
		//   R = renamed
		//   C = copied
		//   U = updated but unmerged
		//   ? = untracked
		//   ! = ignored
		//   T = type change
		//   X = unknown

		status := line[:2]
		path := strings.TrimSpace(line[3:])

		var change string
		switch {
		case strings.HasPrefix(status, "M"):
			change = "modified: " + path
		case strings.HasPrefix(status, "A"):
			change = "added: " + path
		case strings.HasPrefix(status, "D"):
			change = "deleted: " + path
		case strings.HasPrefix(status, "R"):
			change = "renamed: " + path
		case strings.HasPrefix(status, "C"):
			change = "copied: " + path
		case strings.HasPrefix(status, "U"):
			change = "unmerged: " + path
		case strings.HasPrefix(status, "?"):
			change = "untracked: " + path
		case strings.HasPrefix(status, "!"):
			change = "ignored: " + path
		case strings.HasPrefix(status, "T"):
			change = "type changed: " + path
		default:
			change = "unknown: " + path
		}

		changes = append(changes, change)
	}

	return changes
}
