// extractor is used to mine valuable information from a directory and its contents
package extractor

// ExtractedContent represents metadata extracted from a file.
// This is the common format returned by all extractor implementations.
type ExtractedContent struct {
	Category string         // Type of content (e.g., "code", "document", "binary")
	Preview  string         // First few lines of content for quick viewing
	Lines    int            // Total line count (for text files)
	Details  map[string]any // Additional metadata specific to file type
}

// Extractor defines the interface for extracting content from files.
// Each file type (PDF, code, images, etc.) has its own implementation.
type Extractor interface {
	// Extract reads a file and returns structured metadata about its contents.
	//
	// Parameters:
	//   - path: Full path to the file to extract
	//
	// Returns:
	//   - *ExtractedContent: Extracted metadata and preview
	//   - error: Any error encountered during extraction
	Extract(path string) (*ExtractedContent, error)
}
