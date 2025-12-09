package shell

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/DeleMike/scout/internal/scout"
	"github.com/DeleMike/scout/internal/summarize"
)

// HandleScout encapsulates the logic for the "sc" command
func HandleScout(args []string, defaultWriter io.Writer) error {
	var writer io.Writer = defaultWriter
	var targetFile *os.File

	// We filter the args so the rest of the logic doesn't see the ">> file.txt" part
	cleanArgs, fileWriter, err := setupRedirection(args)
	if err != nil {
		return err
	}

	if fileWriter != nil {
		writer = fileWriter
		targetFile = fileWriter
		defer targetFile.Close()
	}

	rawPath := "."
	if len(cleanArgs) > 1 {
		rawPath = strings.Trim(cleanArgs[1], "\"'")
	}

	targetDir, err := filepath.Abs(rawPath)
	if err != nil {
		return fmt.Errorf("error resolving path: %v", err)
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' does not exist", targetDir)
	}

	// Feedback to user (on screen if piping, or to file if not)
	if targetFile != nil {
		fmt.Printf("🔎 Scouting: %s\n", targetDir)
	} else {
		fmt.Fprintf(writer, "🔎 Scouting: %s\n", targetDir)
	}

	summary, insight, err := scout.Run(targetDir)
	if err != nil {
		return err
	}

	fmt.Fprintf(writer, "✅ Found %d files (%.0f%% confidence: %s domain)\n",
		summary.FileCount,
		insight.Confidence*100,
		insight.Domain)

	// Run AI Summarization
	if targetFile != nil {
		fmt.Println("🤖 Generating AI insights...")
	} else {
		fmt.Println("🤖 Generating AI insights...")
	}
	useColor := (targetFile == nil)
	fullPrompt := scout.GeneratePrompt(insight, summary)
	aiResponse, err := summarize.Summarize(fullPrompt, useColor)
	if err != nil {
		return fmt.Errorf("summarizer error: %v", err)
	}

	fmt.Fprintf(writer, "\n%s\n", strings.Repeat("=", 80))
	fmt.Fprintln(writer, aiResponse)
	fmt.Fprintf(writer, "%s\n", strings.Repeat("=", 80))

	if targetFile != nil {
		fmt.Println("✅ Done.")
	}

	return nil
}

// Helper function to extract ">> filename" from args
func setupRedirection(args []string) (cleanArgs []string, file *os.File, err error) {
	redirectIndex := -1
	outputFileName := ""
	isAppend := false

	for i, arg := range args {
		if arg == ">" || arg == ">>" {
			redirectIndex = i
			if arg == ">>" {
				isAppend = true
			}
			if i+1 < len(args) {
				outputFileName = args[i+1]
			}
			break
		}
	}

	if redirectIndex == -1 {
		return args, nil, nil
	}

	if outputFileName == "" {
		return args, nil, fmt.Errorf("no output file specified")
	}

	flags := os.O_CREATE | os.O_WRONLY
	if isAppend {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	f, err := os.OpenFile(outputFileName, flags, 0644)
	if err != nil {
		return args, nil, fmt.Errorf("failed to open file: %w", err)
	}

	fmt.Printf("📝 Saving output to %s...\n", outputFileName)

	// Return args with the redirection part sliced off
	return args[:redirectIndex], f, nil
}
