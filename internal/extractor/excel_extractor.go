package extractor

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExcelExtractor extracts excel file contents
type ExcelExtractor struct{}

// Extract extracts content from an excel file or CSV
func (e ExcelExtractor) Extract(path string) (*ExtractedContent, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get first sheet name
	sheetName := f.GetSheetName(0)

	// Get rows (limit to top 5)
	rows, err := f.GetRows(sheetName, excelize.Options{RawCellValue: true})
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	lineCount := 0
	for _, row := range rows {
		if lineCount > 5 {
			break
		}
		// Join columns with pipe for readability
		sb.WriteString(strings.Join(row, " | "))
		sb.WriteString("\n")
		lineCount++
	}

	return &ExtractedContent{
		Category: "spreadsheet",
		Preview:  sb.String(),
		Details: map[string]any{
			"sheet_name":           sheetName,
			"total_rows_estimated": len(rows),
		},
	}, nil
}
