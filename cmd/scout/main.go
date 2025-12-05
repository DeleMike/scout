package main

import "github.com/DeleMike/scout/internal/shell"

// main is the application entry point.
// It initializes and starts an interactive shell session
// that accepts commands for directory analysis.
func main() {
	// Create a new shell instance with default configuration.
	s := shell.New()

	// Start the REPL (Read-Eval-Print Loop)
	s.Start()
}
