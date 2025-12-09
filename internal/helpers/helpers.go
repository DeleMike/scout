// Package helpers contains various helpful functions
package helpers

import (
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// IsCodeFile checks if a file extension represents source code
func IsCodeFile(ext string) bool {
	// as far as I can add
	codeExts := []string{".go", ".dart", ".js", ".ts", ".py", ".java", ".rb", ".rs", ".c", ".cpp", ".cs", ".php", ".swift", ".kt"}
	return slices.Contains(codeExts, ext)
}

// IsConfigFile checks if a file is a configuration file
// based on extension or filename
func IsConfigFile(ext, name string) bool {
	return ext == ".json" || ext == ".yaml" || ext == ".yml" || ext == ".toml" ||
		ext == ".xml" || ext == ".ini" || name == ".env" || name == ".gitignore"
}

// IsImageFile checks if a file extension represents an image
func IsImageFile(ext string) bool {
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" ||
		ext == ".bmp" || ext == ".svg" || ext == ".webp" || ext == ".heic"
}

// IsVideoFile checks if a file extension represents a video
func IsVideoFile(ext string) bool {
	return ext == ".mp4" || ext == ".avi" || ext == ".mov" || ext == ".mkv" ||
		ext == ".webm" || ext == ".flv" || ext == ".wmv"
}

// IsAudioFile checks if a file extension represents audio
func IsAudioFile(ext string) bool {
	return ext == ".mp3" || ext == ".wav" || ext == ".flac" || ext == ".aac" ||
		ext == ".ogg" || ext == ".m4a" || ext == ".wma"
}

// IsProjectMarkerFile checks if a filename indicates a software project
// (package managers, build files, etc.)
func IsProjectMarkerFile(name string) bool {
	markers := []string{"package.json", "pubspec.yaml", "go.mod", "Cargo.toml",
		"requirements.txt", "pom.xml", "build.gradle", "Gemfile", "composer.json"}
	lowerName := strings.ToLower(name)
	for _, marker := range markers {
		if lowerName == strings.ToLower(marker) {
			return true
		}
	}
	return false
}

// ExtractYear finds a 4-digit year (19xx or 20xx) in a filename
//
// Returns: Year as string, or empty string if not found
func ExtractYear(filename string) string {
	re := regexp.MustCompile(`\b(19|20)\d{2}\b`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

// ExtractTopicsFromFilename extracts meaningful words from a filename
// by removing common separators, numbers, and stop words.
//
// Example: "final_exam_2024_chapter-5.pdf" → ["final", "exam", "chapter"]
//
// Returns: List of topic words (minimum 4 characters, excluding stop words)
func ExtractTopicsFromFilename(filename string) []string {
	// Remove extension and split by common separators
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	name = strings.ToLower(name)

	// Replace separators with spaces
	name = regexp.MustCompile(`[_\-.]`).ReplaceAllString(name, " ")

	// Remove numbers and years
	name = regexp.MustCompile(`\b\d+\b`).ReplaceAllString(name, "")

	// Split into words
	words := strings.Fields(name)

	// Filter out common stop words and short words
	stopWords := map[string]bool{
		"the": true, "and": true, "or": true, "of": true, "a": true, "an": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "with": true,
	}

	topics := []string{}
	for _, word := range words {
		if len(word) > 3 && !stopWords[word] {
			topics = append(topics, word)
		}
	}

	return topics
}

// ShouldPrioritizeDoc determines if a document filename suggests importance
// based on keywords and recency.
//
// Priority keywords: summary, overview, final, important, guide, etc.
// Also prioritizes documents from the last 3 years.
//
// Returns: true if the document should be highlighted
func ShouldPrioritizeDoc(filename string) bool {
	lowerName := strings.ToLower(filename)
	currentYear := time.Now().Year()

	yearStrings := []string{
		strconv.Itoa(currentYear - 2),
		strconv.Itoa(currentYear - 1),
		strconv.Itoa(currentYear),
	}
	priorities := []string{"summary", "overview", "final", "important", "guide",
		"index", "table", "contents", "readme", "start", "intro", "introduction"}

	// Add dynamic year strings
	priorities = append(priorities, yearStrings...)

	for _, priority := range priorities {
		if strings.Contains(lowerName, priority) {
			return true
		}
	}
	return false
}
