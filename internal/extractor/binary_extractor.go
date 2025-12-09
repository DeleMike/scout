package extractor

import "os"

// BinaryExtractor extracts contents from a word document
type BinaryExtractor struct{}

// Extract tried to extract a binary content
func (b BinaryExtractor) Extract(path string) (*ExtractedContent, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &ExtractedContent{
		Category: "binary",
		Preview:  "",
		Lines:    0,
		Details: map[string]any{
			"size_bytes": info.Size(),
		},
	}, nil
}
