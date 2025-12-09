package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DeleMike/scout/internal/scout"
	"github.com/DeleMike/scout/internal/summarize"
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
	case "sc":
		rawPath := "."
		if len(args) > 1 {
			rawPath = strings.Trim(args[1], "\"'")
		}

		targetDir, err := filepath.Abs(rawPath)
		if err != nil {
			fmt.Printf("❌ Error resolving path: %v\n", err)
			return true
		}

		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			fmt.Printf("❌ Error: Directory '%s' does not exist.\n", targetDir)
			return true
		}

		fmt.Printf("🔎 Scouting: %s\n", targetDir)

		summary, insight, err := scout.Run(targetDir)
		if err != nil {
			fmt.Println("Error:", err)
			return true
		}

		fmt.Printf("✅ Found %d files (%.0f%% confidence: %s domain)\n",
			summary.FileCount,
			insight.Confidence*100,
			insight.Domain)

		fmt.Println("🤖 Generating AI insights...")
		fullPrompt := scout.GeneratePrompt(insight, summary)

		aiResponse, err := summarize.Summarize(fullPrompt)
		if err != nil {
			fmt.Println("Summarizer error:", err)
			return true
		}

		fmt.Println("\n" + strings.Repeat("=", 120))
		fmt.Println(aiResponse)
		fmt.Println(strings.Repeat("=", 120))
		return true

	}
	return false
}
