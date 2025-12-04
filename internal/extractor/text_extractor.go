package extractor

import (
	"io"
	"os"
	"slices"
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
		".c", ".cpp", ".h", ".hpp", ".cmake", ".sh",
	}

	return slices.Contains(textExt, ext)
}

type GenericTextExtractor struct{}

func (e GenericTextExtractor) Extract(path string) (*ExtractedContent, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read ONLY the first 1000 bytes
	buf := make([]byte, 1000)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	text := string(buf[:n])

	lines := 0
	for _, c := range text {
		if c == '\n' {
			lines++
		}
	}

	return &ExtractedContent{
		Category: "text",
		Preview:  text,
		Lines:    lines,
		Details:  map[string]any{"format": "text"},
	}, nil
}
