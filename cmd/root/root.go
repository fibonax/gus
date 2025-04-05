package root

import (
	"github.com/nguyendangminh/gus/pkg/core"
	"github.com/spf13/cobra"
)

var (
	// jsonOutput determines if the output should be in JSON format
	jsonOutput bool
	// rootPath is the path to scan for Git repositories
	rootPath string
	// verbose determines if verbose output should be shown
	verbose bool
)

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gus",
		Short: "Git Uncommitted Scanner - Find Git repositories with uncommitted changes",
		Long: `A command-line tool to scan directories for Git repositories with uncommitted changes.
It recursively searches through directories to find Git repositories and checks their status.`,
		RunE: run,
	}

	// Add flags
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	cmd.Flags().StringVar(&rootPath, "path", ".", "path to scan for Git repositories")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	return cmd
}

// run is the main function that will be executed when the command is run
func run(cmd *cobra.Command, args []string) error {
	// If path is provided as an argument, override the flag
	if len(args) > 0 {
		rootPath = args[0]
	}

	// Create scanner with options
	options := core.Options{
		Path:    rootPath,
		JSON:    jsonOutput,
		Verbose: verbose,
	}
	scanner := core.New(options)

	// Run the scanner
	return scanner.Run()
}
