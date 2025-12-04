package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type FileType int

const (
	File FileType = iota
	Directory
)

type FileInfo struct {
	Name    string
	Path    string
	Type    FileType
	FileExt string
	Size    int64
}

type ScanResult struct {
	Path           string
	Files          []FileInfo
	Subdirectories []string
}

// ScanDirectory checks what are the contents of the current working directory
func ScanDirectory(root string) (*ScanResult, error) {
	summary := &ScanResult{
		Path: root,
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		name := d.Name()

		if strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() && path != root {
			summary.Subdirectories = append(summary.Subdirectories, path)
			return nil
		}

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
