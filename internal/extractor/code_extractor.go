package extractor

import (
	"os"
	"strings"
)

// CodeExtractor extracts contents from a code file
type CodeExtractor struct{}

// Extract extracts content from a code file
func (c CodeExtractor) Extract(path string) (*ExtractedContent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	preview := strings.Join(lines[:min(len(lines), 30)], "\n")

	return &ExtractedContent{
		Category: "code",
		Preview:  preview,
		Lines:    len(lines),
		Details: map[string]any{
			"imports": extractImports(lines),
		},
	}, nil
}

func extractImports(lines []string) []string {
	var imports []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "import ") {
			imports = append(imports, l)
		}
	}
	return imports
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
