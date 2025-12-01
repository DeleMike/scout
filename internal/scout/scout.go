package scout

import (
	"fmt"

	"github.com/DeleMike/scout/internal/extractor"
	"github.com/DeleMike/scout/internal/scanner"
)

// FileSummary is the structured info for each file
type FileSummary struct {
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Size     int64          `json:"size_bytes"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// DirectorySummary is the structured summary of a directory
type DirectorySummary struct {
	Directory      string        `json:"directory"`
	FileCount      int           `json:"file_count"`
	Subdirectories []string      `json:"subdirectories"`
	Files          []FileSummary `json:"files"`
}

// Run scans a directory, extracts content metadata, and returns structured summary
func Run(root string) (*DirectorySummary, error) {
	dir, err := scanner.ScanDirectory(root)
	if err != nil {
		return nil, err
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

		extractor := extractor.DetectCategory(file.FileExt)
		content, err := extractor.Extract(file.Path)

		fileSummary := FileSummary{
			Name: file.Name,
			Type: "unknown",
			Size: file.Size,
		}

		if err != nil {
			fmt.Printf("[%s] extraction error: %v\n", file.Name, err)
		} else {
			fileSummary.Type = content.Category
			fileSummary.Metadata = map[string]any{
				"preview": content.Preview,
				"lines":   content.Lines,
			}
		}

		summary.Files = append(summary.Files, fileSummary)

	}
	return summary, nil

}
