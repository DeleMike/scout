package shell

import (
	"fmt"
	"os"
	"strings"
)

// runBuiltin executes built-in shell commands that don't require
// external processes.
//
// Supported commands:
//   - exit: Terminate the shell
//   - pwd: Print working directory
//   - ls: List directory contents
//   - sc: Run Scout directory analysis
//
// Parameters:
//   - args: Command and its arguments (args[0] is the command name)
//
// Returns:
//   - bool: true if command was recognized and handled, false otherwise.
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
			if strings.HasPrefix(ent.Name(), ".") {
				continue
			}
			fmt.Println(ent.Name())
		}
		return true
	case "scout", "sc":
		err := HandleScout(args, os.Stdout)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
		}
		return true
	}
	return false
}
