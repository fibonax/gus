package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Exit with the same code
	os.Exit(code)
}
