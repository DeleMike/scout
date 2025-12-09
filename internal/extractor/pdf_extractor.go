package extractor

import (
	"bytes"

	"github.com/ledongthuc/pdf"
)

// PDFExtractor extracts content from a pdf file
type PDFExtractor struct{}

// Extract extracts content from a pdf file
func (e PDFExtractor) Extract(path string) (*ExtractedContent, error) {
	f, r, err := pdf.Open(path)
	if err != nil {

		return &ExtractedContent{
			Category: "document",
			Preview:  "[PDF content could not be extracted - likely encrypted or unsupported format]",
			Details: map[string]any{
				"type":  "pdf",
				"error": "read_failed",
			},
		}, nil
	}
	defer f.Close()

	var buf bytes.Buffer
	limit := 3
	if r.NumPage() < limit {
		limit = r.NumPage()
	}

	for i := 1; i <= limit; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		text, _ := p.GetPlainText(nil)
		buf.WriteString(text)
		buf.WriteString("\n")
	}

	preview := buf.String()

	// Safety Truncation
	if len(preview) > 1000 {
		preview = preview[:1000] + "..."
	}

	if len(preview) == 0 {
		preview = "[Scanned PDF or Image-based - No text extracted]"
	}

	return &ExtractedContent{
		Category: "document",
		Preview:  preview,
		Details: map[string]any{
			"pages": r.NumPage(),
			"type":  "pdf",
		},
	}, nil
}
