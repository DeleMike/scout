package extractor

import (
	"slices"
	"os"
	"strings"
)

// IsTextFile checks by extension
func IsTextFile(ext string) bool {
	ext = strings.ToLower(ext)
	textExt := []string{
		".txt", ".md", ".json", ".yaml", ".yml",
		".pdf", ".doc", "docx",
		".xml", ".csv", ".ini", ".cfg",
		".go", ".py", ".js", ".ts", ".java", ".swift",
		".rb", ".rs", ".php", ".css", ".html",
	}

	return slices.Contains(textExt, ext)
}

type GenericTextExtractor struct{}

func (g GenericTextExtractor) Extract(path string) (*ExtractedContent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := string(data)
	lines := strings.Split(raw, "\n")

	preview := strings.Join(lines[:min(len(lines), 20)], "\n")

	return &ExtractedContent{
		Category: "text",
		Preview:  preview,
		Lines:    len(lines),
		Details:  map[string]any{},
	}, nil
}
