package main

import "github.com/DeleMike/scout/internal/shell"

func main() {
	// create new shell and start
	s := shell.New()
	s.Start()
}
