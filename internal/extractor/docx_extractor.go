package extractor

import (
	"archive/zip"
	"encoding/xml"
	"strings"
)

type DocxExtractor struct{}

func (e DocxExtractor) Extract(path string) (*ExtractedContent, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var content strings.Builder

	// Find the main content XML inside the zip
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}

			// Simple XML parsing to just grab text content
			decoder := xml.NewDecoder(rc)
			for {
				t, _ := decoder.Token()
				if t == nil {
					break
				}
				switch se := t.(type) {
				case xml.CharData:
					content.Write(se)
					content.WriteString(" ")
				}
			}
			rc.Close()
			break
		}
	}

	text := content.String()
	// Clean up XML noise if necessary, or just truncate
	if len(text) > 1000 {
		text = text[:1000] + "..."
	}

	return &ExtractedContent{
		Category: "document",
		Preview:  text,
		Details: map[string]any{
			"type": "docx",
		},
	}, nil
}
