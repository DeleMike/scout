package shell

import (
	"fmt"
	"os"

	"github.com/DeleMike/scout/internal/scanner"
)

// runBuiltin runs simple commands that every shell application should have.
//
// It is a typically only runs one command. (e.g, cd, ls, pwd, etc). However, it could be more complex.
// Basically, it consists of the usual thing you think a typical Unix terminal can run.
func (s *Shell) runBuiltin(args []string) bool {
	switch args[0] {
	case "exit":
		fmt.Print("Bye, scout!")
		os.Exit(0)
	case "pwd":
		wd, _ := os.Getwd()
		fmt.Println(wd)
		return true
	case "ls":
		entries, _ := os.ReadDir(".")
		for _, ent := range entries {
			fmt.Println(ent.Name())
		}
		return true
	case "scout":
		wd, _ := os.Getwd()
		scannedDir, err := scanner.ScanDirectory(wd)
		if err != nil {
			fmt.Println("Something happened")
			return true
		}

		fmt.Printf("%v\n", scannedDir.Pretty())

		return true

	}
	return false
}
