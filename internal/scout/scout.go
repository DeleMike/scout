package scout

import (
	"fmt"

	"github.com/DeleMike/scout/internal/extractor"
	"github.com/DeleMike/scout/internal/scanner"
)

// FileSummary represents structured metadata for a single file
type FileSummary struct {
	Name      string         `json:"name"`               // Filename
	Type      string         `json:"type"`               // Category (code, document, etc.)
	Extension string         `json:"extension"`          // File extension
	Size      int64          `json:"size_bytes"`         // Size in bytes
	Metadata  map[string]any `json:"metadata,omitempty"` // Extracted content details
}

// DirectorySummary is the complete analysis result for a directory
type DirectorySummary struct {
	Directory      string        `json:"directory"`      // Root path
	FileCount      int           `json:"file_count"`     // Total files
	Subdirectories []string      `json:"subdirectories"` // Subdirectory paths
	Files          []FileSummary `json:"files"`          // List of files in Directory
}

// Run is the main entry point for directory analysis.
// It orchestrates the entire pipeline:
//  1. Scan directory structure
//  2. Extract content from each file
//  3. Analyze patterns and generate insights
//
// Parameters:
//   - root: Path to directory to analyze
//
// Returns:
//   - *DirectorySummary: Structured file metadata
//   - *ContentInsight: AI-ready analysis and recommendations
//   - error: Any error encountered during processing
func Run(root string) (*DirectorySummary, *ContentInsight, error) {
	// Scan directory structure
	dir, err := scanner.ScanDirectory(root)
	if err != nil {
		return nil, nil, err
	}

	summary := &DirectorySummary{
		Directory:      root,
		Subdirectories: dir.Subdirectories,
		FileCount:      len(dir.Files),
	}

	for _, file := range dir.Files {
		if file.Type == scanner.Directory {
			continue
		}

		// Get appropriate extractor for this file type
		extractor := extractor.DetectCategory(file.FileExt)
		content, err := extractor.Extract(file.Path)

		fileSummary := FileSummary{
			Name:      file.Name,
			Type:      "unknown",
			Extension: file.FileExt,
			Size:      file.Size,
		}

		if err != nil {
			fmt.Printf("[%s] extraction error: %v\n", file.Name, err)
		} else {
			fileSummary.Type = content.Category
			fileSummary.Metadata = map[string]any{
				"preview": content.Preview,
				"lines":   content.Lines,
				"details": content.Details,
			}
		}

		summary.Files = append(summary.Files, fileSummary)

	}

	// Analyze directory to generate insights
	insight := AnalyzeDirectory(summary)

	return summary, insight, nil

}
