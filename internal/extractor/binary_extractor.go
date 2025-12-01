package extractor

import "os"

type BinaryExtractor struct{}

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
