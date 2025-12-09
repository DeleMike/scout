// Package scanner is used to read a directory to know what its contents entail
package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// FileType represents the type of file system entry.
type FileType int

const (
	// File is a Regular file
	File FileType = iota
	// Directory is a Directory/folder
	Directory
)

// FileInfo contains metadata about a single file.
type FileInfo struct {
	Name    string   // Base filename (e.g., "main.go")
	Path    string   // Full path to file
	Type    FileType // File or Directory
	FileExt string   // File extension (e.g., ".go")
	Size    int64    // Size in bytes
}

// ScanResult contains the complete scan of a directory.
type ScanResult struct {
	Path           string     // Root directory path
	Files          []FileInfo // All files found
	Subdirectories []string   // Paths to subdirectories
}

// ScanDirectory recursively walks a directory tree and collects
// information about all files and subdirectories.
//
// Hidden files (starting with '.') are automatically excluded to
// avoid scanning system files, git directories, etc.
//
// Parameters:
//   - root: Path to directory to scan
//
// Returns:
//   - *ScanResult: Complete directory structure
//   - error: Any error encountered during scanning
func ScanDirectory(root string) (*ScanResult, error) {
	summary := &ScanResult{
		Path: root,
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip entries that can't be read
		}

		name := d.Name()

		// skip hidden files
		if strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir // Don't recurse into hidden dirs
			}
			return nil
		}

		// Track subdirectories (but don't add root itself)
		if d.IsDir() && path != root {
			summary.Subdirectories = append(summary.Subdirectories, path)
			return nil
		}

		// Add regular files to the result
		if !d.IsDir() {
			info, _ := d.Info()
			fileExt := strings.ToLower(filepath.Ext(name))

			summary.Files = append(summary.Files, FileInfo{
				Name:    name,
				Path:    path,
				Type:    File,
				FileExt: fileExt,
				Size:    info.Size(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return summary, nil
}

// Pretty formats the scan result as a human-readable string
// showing files with their extensions and sizes, plus subdirectories.
//
// Returns:
//   - string: Formatted directory listing
func (d *ScanResult) Pretty() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Directory: %s\n\n", d.Path))
	b.WriteString("Files:\n")

	for _, f := range d.Files {
		if f.Type == Directory {
			b.WriteString(fmt.Sprintf("  - %s (dir)\n", f.Name))
			continue
		}

		b.WriteString(fmt.Sprintf(
			"  - %s (%s, %d bytes)\n",
			f.Name,
			f.FileExt,
			f.Size,
		))
	}

	b.WriteString("\nSubdirectories:\n")
	for _, sub := range d.Subdirectories {
		rel, _ := filepath.Rel(d.Path, sub)
		b.WriteString(fmt.Sprintf("  - %s\n", rel))
	}

	return b.String()
}
