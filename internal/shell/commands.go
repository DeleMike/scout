package shell

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/DeleMike/scout/internal/scout"
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
	case "sc":
		wd, _ := os.Getwd()
		summary, err := scout.Run(wd)
		if err != nil {
			fmt.Println("Error:", err)
			return true
		}

		prettyJSON, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			fmt.Println("JSON marshal error:", err)
			return true
		}
		fmt.Println(string(prettyJSON))
		return true

	}
	return false
}
