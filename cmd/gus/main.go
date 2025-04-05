package main

import (
	"fmt"
	"os"

	"github.com/nguyendangminh/gus/cmd/root"
)

func main() {
	if err := root.NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// TODO: Implement main logic
	return nil
}
