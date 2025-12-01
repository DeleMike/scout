package extractor

import (
	"os"
	"strings"
)

type MarkdownExtractor struct{}

func (m MarkdownExtractor) Extract(path string) (*ExtractedContent, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	title := ""
	for _, l := range lines {
		if after, ok :=strings.CutPrefix(l, "# "); ok  {
			title = after
			break
		}
	}

	preview := strings.Join(lines[:min(len(lines), 20)], "\n")

	return &ExtractedContent{
		Category: "markdown",
		Preview:  preview,
		Lines:    len(lines),
		Details: map[string]any{
			"title": title,
		},
	}, nil
}
