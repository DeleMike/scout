package extractor

// DetectCategory determines the appropriate extractor for a file
// based on its extension.
//
// This acts as a factory pattern, routing files to specialized
// extractors that understand their format.
//
// Parameters:
//   - ext: File extension including the dot (e.g., ".go", ".pdf")
//
// Returns:
//   - Extractor: Appropriate extractor implementation for the file type
//
// File categories:
//   - Code: .go, .dart, .js, .ts, .py, .java, .rb, .rs, .c, .cpp
//   - PDF: .pdf
//   - Word: .docx, .doc
//   - Excel: .xlsx, .xls
//   - Text: .md, .txt
//   - Structured: .json, .yaml, .xml, .csv, etc.
//   - Binary: Images, audio, video, or unknown formats
func DetectCategory(ext string) Extractor {
	switch ext {
	case ".go", ".dart", ".js", ".ts", ".py", ".java", ".rb", ".rs", ".c", ".cpp":
		return CodeExtractor{}
	case ".pdf":
		return PDFExtractor{}
	case ".docx", ".doc":
		return DocxExtractor{}
	case ".xlsx", ".xls":
		return ExcelExtractor{}
	case ".md", ".txt":
		return MarkdownExtractor{}
	case ".json", ".yaml", ".yml", ".toml", ".env", ".xml", ".csv", ".cmake":
		return GenericTextExtractor{}
	case ".png", ".jpg", ".jpeg", ".mp3", ".mp4":
		return BinaryExtractor{}
	default:
		// Fallback: Check if it's a text file by content
		if IsTextFile(ext) {
			return GenericTextExtractor{}
		}
		return BinaryExtractor{}
	}
}
