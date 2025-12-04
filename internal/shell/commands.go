package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/DeleMike/scout/internal/scanner"
	"github.com/DeleMike/scout/internal/scout"
	"github.com/DeleMike/scout/internal/summarize"
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
		// wd, _ := os.Getwd()
		summary, insight, err := scout.Run("/Users/mac/Documents/Life Changing Books")
		if err != nil {
			fmt.Println("Error:", err)
			return true
		}

		fmt.Printf("✅ Found %d files (%.0f%% confidence: %s domain)\n",
			summary.FileCount,
			insight.Confidence*100,
			insight.Domain)

		// Step 2: Convert scanner results for GeneratePrompt
		// We need to create a ScanResult with the basic file info
		scanResult := &scanner.ScanResult{
			Path:           summary.Directory,
			Subdirectories: summary.Subdirectories,
			Files:          make([]scanner.FileInfo, len(summary.Files)),
		}

		for i, file := range summary.Files {
			scanResult.Files[i] = scanner.FileInfo{
				Name:    file.Name,
				FileExt: file.Extension,
				Size:    file.Size,
				Type:    scanner.File,
			}
		}

		fmt.Println("🤖 Generating AI insights...")
		fullPrompt := scout.GeneratePrompt(insight, scanResult)

		aiResponse, err := summarize.Summarize(fullPrompt)
		if err != nil {
			fmt.Println("Summarizer error:", err)
			return true
		}

		// prettyJSON, err := json.MarshalIndent(summary, "", "  ")
		// if err != nil {
		// 	fmt.Println("JSON marshal error:", err)
		// 	return true
		// }

		// fmt.Printf("scanning result: \n%v",string(prettyJSON))

		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println(aiResponse)
		fmt.Println(strings.Repeat("=", 60))
		return true

	}
	return false
}
