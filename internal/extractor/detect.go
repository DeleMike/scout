package extractor

func DetectCategory(ext string) Extractor {
	switch ext {
	case ".go", ".dart", ".js", ".ts", ".py", ".java", ".rb", ".rs", ".c", ".cpp":
		return CodeExtractor{}
	case ".md", ".txt":
		return MarkdownExtractor{}
	case ".json", ".yaml", ".yml", ".toml", ".env":
		return GenericTextExtractor{}
	case ".png", ".jpg", ".jpeg", ".mp3", ".mp4":
		return BinaryExtractor{}
	default:
		if IsTextFile(ext) {
			return GenericTextExtractor{}
		}
		return BinaryExtractor{}
	}
}
