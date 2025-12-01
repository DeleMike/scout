// extractor is used to mine valuable information from a directory and its contents
package extractor

// ExtractedContent is just an abstraction for a peek into a file (if possible)
type ExtractedContent struct {
	Category string
	Preview  string // just a few lines
	Lines    int
	Details  map[string]any
}

// / Extractor methods
type Extractor interface {
	Extract(path string) (*ExtractedContent, error)
}
