package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nguyendangminh/gus/pkg/git"
)

// FormatOptions contains options for formatting output
type FormatOptions struct {
	JSON bool
}

// FormatRepositories formats the list of repositories according to the options
func FormatRepositories(repos []*git.Repository, opts FormatOptions) error {
	if len(repos) == 0 {
		fmt.Println("No Git repositories with uncommitted changes found.")
		return nil
	}

	if opts.JSON {
		return formatJSON(repos)
	}

	return formatText(repos)
}

// formatJSON formats the repositories as JSON
func formatJSON(repos []*git.Repository) error {
	type repoJSON struct {
		Path       string    `json:"path"`
		Changes    []string  `json:"changes"`
		ScanTime   time.Time `json:"scan_time"`
		TotalRepos int       `json:"total_repositories"`
	}

	jsonRepos := make([]repoJSON, len(repos))
	for i, repo := range repos {
		jsonRepos[i] = repoJSON{
			Path:       repo.Path,
			Changes:    repo.Changes,
			ScanTime:   time.Now(),
			TotalRepos: len(repos),
		}
	}

	// Create a wrapper structure for the entire output
	type outputJSON struct {
		Repositories []repoJSON `json:"repositories"`
		Metadata     struct {
			ScanTime   time.Time `json:"scan_time"`
			TotalRepos int       `json:"total_repositories"`
			Version    string    `json:"version"`
		} `json:"metadata"`
	}

	output := outputJSON{
		Repositories: jsonRepos,
		Metadata: struct {
			ScanTime   time.Time `json:"scan_time"`
			TotalRepos int       `json:"total_repositories"`
			Version    string    `json:"version"`
		}{
			ScanTime:   time.Now(),
			TotalRepos: len(repos),
			Version:    "1.0.0",
		},
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// formatText formats the repositories as text
func formatText(repos []*git.Repository) error {
	fmt.Printf("Found %d Git repositories with uncommitted changes:\n\n", len(repos))

	for i, repo := range repos {
		// Format repository path
		path := repo.Path
		if strings.HasPrefix(path, os.Getenv("HOME")) {
			path = "~" + path[len(os.Getenv("HOME")):]
		}

		fmt.Printf("%d. %s\n", i+1, path)

		// Format changes
		for _, change := range repo.Changes {
			fmt.Printf("   - %s\n", change)
		}
		fmt.Println()
	}

	return nil
}
