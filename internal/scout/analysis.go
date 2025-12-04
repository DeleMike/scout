package scout

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/DeleMike/scout/internal/utils"
)

// DomainType represents the detected purpose of a directory
type DomainType string

// during analysis, we try to group the different preview/details into groups for better insights on what scout is looking at
const (
	DomainSoftwareProject DomainType = "software"
	DomainDocuments       DomainType = "documents"
	DomainMedia           DomainType = "media"
	DomainStudyMaterials  DomainType = "study"
	DomainFinancial       DomainType = "financial"
	DomainCreative        DomainType = "creative"
	DomainMixed           DomainType = "mixed"
	DomainEmpty           DomainType = "empty"
)

type ContentInsight struct {
	Domain          DomainType
	Topics          []string       // Main theme
	DateRange       string         // for time-sensitive content
	KeyFiles        []string       // most important files(presumably)
	FilesByCategory map[string]int // domain type of files in directory
	Recommendations []string       // our guess
	Confidence      float64        // confidence of our guess
}

func AnalyzeDirectory(summary *DirectorySummary) *ContentInsight {
	insight := &ContentInsight{
		Topics:          []string{},
		KeyFiles:        []string{},
		FilesByCategory: make(map[string]int),
		Recommendations: []string{},
	}

	categories := categorizeFiles(summary.Files)
	insight.FilesByCategory = categories

	// Detect domain based on file distribution and content
	insight.Domain = detectDomain(categories, summary.Files)
	insight.Confidence = calculateConfidence(categories)

	// Extract insights based on domain
	switch insight.Domain {
	case DomainSoftwareProject:
		extractSoftwareInsights(insight, summary.Files)
	case DomainDocuments, DomainStudyMaterials:
		extractDocumentInsights(insight, summary.Files)
	case DomainMedia:
		extractMediaInsights(insight, summary.Files)
	case DomainFinancial:
		extractFinancialInsights(insight, summary.Files)
	case DomainCreative:
		extractCreativeInsights(insight)
	default:
		extractMixedInsights(insight, summary.Files)
	}

	return insight
}

func categorizeFiles(files []FileSummary) map[string]int {
	categories := make(map[string]int)

	for _, file := range files {
		ext := strings.ToLower(file.Extension)
		name := strings.ToLower(file.Name)

		switch {
		// software
		case utils.IsCodeFile(ext):
			categories["code"]++
		case utils.IsConfigFile(ext, name):
			categories["config"]++
			// Documents
		case ext == ".pdf":
			categories["pdf"]++
		case ext == ".docx" || ext == ".doc":
			categories["word"]++
		case ext == ".xlsx" || ext == ".xls" || ext == ".csv":
			categories["spreadsheet"]++
		case ext == ".pptx" || ext == ".ppt":
			categories["presentation"]++
		case ext == ".txt" || ext == ".md":
			categories["text"]++

		// Media
		case utils.IsImageFile(ext):
			categories["image"]++
		case utils.IsVideoFile(ext):
			categories["video"]++
		case utils.IsAudioFile(ext):
			categories["audio"]++

		// Archives
		case ext == ".zip" || ext == ".rar" || ext == ".7z" || ext == ".tar" || ext == ".gz":
			categories["archive"]++

		default:
			categories["other"]++
		}
	}

	return categories

}

func detectDomain(categories map[string]int, files []FileSummary) DomainType {
	total := 0
	for _, count := range categories {
		total += count
	}

	if total == 0 {
		return DomainEmpty
	}

	// Calculate percentages of different type of domains available (% dist)
	codePercent := float64(categories["code"]+categories["config"]) / float64(total)
	docPercent := float64(categories["pdf"]+categories["word"]+categories["text"]) / float64(total)
	mediaPercent := float64(categories["image"]+categories["video"]+categories["audio"]) / float64(total)
	spreadsheetPercent := float64(categories["spreadsheet"]) / float64(total)

	// Check for software project markers (mostly a strong indication that the directory is a software)
	hasProjectMarkers := false
	for _, file := range files {
		if utils.IsProjectMarkerFile(file.Name) {
			hasProjectMarkers = true
			break
		}
	}

	// Domain detection logic
	if hasProjectMarkers || codePercent > 0.3 {
		return DomainSoftwareProject
	}

	if docPercent > 0.5 {
		// Check if study materials
		if hasStudyKeywords(files) {
			return DomainStudyMaterials
		}
		// Check if financial docs
		if hasFinancialKeywords(files) {
			return DomainFinancial
		}
		return DomainDocuments
	}

	if mediaPercent > 0.7 {
		// Mostly images with creative file names
		if categories["image"] > categories["video"]+categories["audio"] && hasCreativeKeywords(files) {
			return DomainCreative
		}
		return DomainMedia
	}

	if spreadsheetPercent > 0.4 {
		return DomainFinancial
	}

	return DomainMixed

}

// calculateConfidence measures how dominant the most frequent category is relative to all categories.
func calculateConfidence(categories map[string]int) float64 {
	total := 0
	dominant := 0

	for _, count := range categories {
		total += count
		if count > dominant {
			dominant = count
		}
	}

	if total == 0 {
		return 0
	}

	return float64(dominant) / float64(total)
}

func hasStudyKeywords(files []FileSummary) bool {
	keywords := []string{"exam", "test", "quiz", "study", "lecture", "notes", "chapter",
		"assignment", "homework", "course", "syllabus", "practice", "review"}

	matchCount := 0
	for _, file := range files {
		lowerName := strings.ToLower(file.Name)
		for _, keyword := range keywords {
			if strings.Contains(lowerName, keyword) {
				matchCount++
				break
			}
		}
	}
	return matchCount >= 3
}

func hasFinancialKeywords(files []FileSummary) bool {
	keywords := []string{"tax", "invoice", "receipt", "statement", "bank", "payroll",
		"expense", "budget", "financial", "accounting"}

	matchCount := 0
	for _, file := range files {
		lowerName := strings.ToLower(file.Name)
		for _, keyword := range keywords {
			if strings.Contains(lowerName, keyword) {
				matchCount++
				break
			}
		}
	}
	return matchCount >= 2
}

func hasCreativeKeywords(files []FileSummary) bool {
	keywords := []string{"design", "mockup", "draft", "sketch", "artwork", "render",
		"illustration", "logo", "banner", "poster"}

	matchCount := 0
	for _, file := range files {
		lowerName := strings.ToLower(file.Name)
		for _, keyword := range keywords {
			if strings.Contains(lowerName, keyword) {
				matchCount++
				break
			}
		}
	}
	return matchCount >= 2
}

// Domain-specific insight extraction
func extractDocumentInsights(insight *ContentInsight, files []FileSummary) {
	// Analyze PDF filenames for topics and dates
	topics := make(map[string]int)
	years := make(map[string]bool)

	for _, file := range files {
		if file.Extension == ".pdf" || file.Extension == ".docx" {

			// extract year if present
			if year := utils.ExtractYear(file.Name); year != "" {
				years[year] = true
			}

			// extract potential topics from filename
			fileTopics := utils.ExtractTopicsFromFilename(file.Name)
			for _, topic := range fileTopics {
				topics[topic]++
			}

			// Prioritize files for key files list
			if utils.ShouldPrioritizeDoc(file.Name) {
				insight.KeyFiles = append(insight.KeyFiles, file.Name)
			}

		}
	}
	for topic := range topics {
		insight.Topics = append(insight.Topics, topic)
	}

	// Set date range if years found
	if len(years) > 0 {
		minYear, maxYear := "", ""
		for year := range years {
			if minYear == "" || year < minYear {
				minYear = year
			}
			if maxYear == "" || year > maxYear {
				maxYear = year
			}
		}
		if minYear == maxYear {
			insight.DateRange = minYear
		} else {
			insight.DateRange = fmt.Sprintf("%s-%s", minYear, maxYear)
		}
	}

	// Generate recommendations
	if len(insight.KeyFiles) == 0 {
		// Find most recent or important looking files
		insight.KeyFiles = findImportantDocs(files, 3)
	} else if len(insight.KeyFiles) > 5 {
		insight.KeyFiles = insight.KeyFiles[:5]
	}

	insight.Recommendations = []string{
		"Start with the most recent documents",
		"Look for summary or overview files first",
		"Organize by topic or date if needed",
	}
}

func extractSoftwareInsights(insight *ContentInsight, files []FileSummary) {
	// Detect tech stack
	stacks := detectTechStack(files)
	insight.Topics = stacks

	// Find entry points
	entryPoints := []string{"main.dart", "main.go", "index.js", "index.html", "app.py", "main.py"}
	for _, file := range files {
		for _, entry := range entryPoints {
			if strings.ToLower(file.Name) == entry {
				insight.KeyFiles = append(insight.KeyFiles, file.Name)
			}
		}
	}

	// Find README
	for _, file := range files {
		if strings.ToLower(file.Name) == "readme.md" {
			insight.KeyFiles = append([]string{file.Name}, insight.KeyFiles...)
			break
		}
	}

	if len(insight.KeyFiles) == 0 {
		insight.KeyFiles = []string{"Look in the src/ or lib/ directory"}
	}

	insight.Recommendations = []string{
		"Read README.md first if available",
		"Check the main entry point to understand flow",
		"Review package/dependency files for tech stack",
	}
}

func extractMediaInsights(insight *ContentInsight, files []FileSummary) {
	imageCount := 0
	videoCount := 0
	audioCount := 0

	for _, file := range files {
		ext := strings.ToLower(file.Extension)
		if utils.IsImageFile(ext) {
			imageCount++
		} else if utils.IsVideoFile(ext) {
			videoCount++
		} else if utils.IsAudioFile(ext) {
			audioCount++
		}
	}

	if imageCount > videoCount && imageCount > audioCount {
		insight.Topics = []string{"photos", "images"}
		insight.Recommendations = []string{"Browse through and enjoy the memories! 📸"}
	} else if videoCount > 0 {
		insight.Topics = []string{"videos"}
		insight.Recommendations = []string{"Grab some popcorn and enjoy! 🍿"}
	} else if audioCount > 0 {
		insight.Topics = []string{"music", "audio"}
		insight.Recommendations = []string{"Put on your headphones and vibe! 🎵"}
	}
}

func extractFinancialInsights(insight *ContentInsight, files []FileSummary) {
	topics := make(map[string]bool)

	for _, file := range files {
		name := strings.ToLower(file.Name)
		if strings.Contains(name, "tax") {
			topics["taxes"] = true
		}
		if strings.Contains(name, "assessment") {
			topics["taxes"] = true
		}
		if strings.Contains(name, "invoice") || strings.Contains(name, "receipt") {
			topics["invoices/receipts"] = true
		}
		if strings.Contains(name, "statement") || strings.Contains(name, "bank") {
			topics["bank statements"] = true
		}
		if strings.Contains(name, "payroll") || strings.Contains(name, "salary") {
			topics["payroll"] = true
		}
	}

	for topic := range topics {
		insight.Topics = append(insight.Topics, topic)
	}

	insight.Recommendations = []string{
		"Organize by year and category",
		"Keep tax documents separate and secure",
		"Back up important financial records",
	}
}

func extractCreativeInsights(insight *ContentInsight) {
	insight.Topics = []string{"creative work", "design assets"}
	insight.Recommendations = []string{
		"Browse through for inspiration",
		"Consider organizing by project or date",
		"Keep high-res originals backed up",
	}
}

func extractMixedInsights(insight *ContentInsight, files []FileSummary) {
	insight.Topics = []string{"various files"}

	type fileScore struct {
		name string
		size int64
	}
	var sorted []fileScore
	for _, f := range files {
		sorted = append(sorted, fileScore{f.Name, f.Size})
	}

	// Sort by size descending
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].size > sorted[j].size
	})

	// Take top 3
	limit := min(len(sorted), 3)

	for i := range limit {
		insight.KeyFiles = append(insight.KeyFiles, sorted[i].name)
	}

	insight.Recommendations = []string{
		"Check the largest files first",
		"Sort by file type to organize",
	}
}

// findImportantDocs sorts files based on a scoring system and returns the top K files,
// or all files if K exceeds the number of files. It includes a fallback mechanism if
// no relevant documents are found.
func findImportantDocs(files []FileSummary, limit int) []string {
	type scoredFile struct {
		name  string
		score int
	}

	var scored []scoredFile

	currentYear := time.Now().Year()
	years := []string{
		strconv.Itoa(currentYear - 2),
		strconv.Itoa(currentYear - 1),
		strconv.Itoa(currentYear),
	}

	for _, file := range files {
		// Filter: Only look at Documents
		ext := strings.ToLower(file.Extension)
		if ext != ".pdf" && ext != ".docx" && ext != ".doc" {
			continue
		}

		score := 0
		name := strings.ToLower(file.Name)

		// 1. Keyword Scoring
		if strings.Contains(name, "summary") || strings.Contains(name, "overview") {
			score += 50
		}
		if strings.Contains(name, "final") || strings.Contains(name, "important") {
			score += 40
		}
		if strings.Contains(name, "guide") || strings.Contains(name, "handbook") {
			score += 30
		}

		// 2. Year Relevance (Recent = Better)
		if strings.Contains(name, years[2]) { // Current year
			score += 30
		} else if strings.Contains(name, years[1]) { // Last year
			score += 20
		} else if strings.Contains(name, years[0]) { // 2 years ago
			score += 10
		}

		// 3. Size Weighting (Larger files often contain the "Core" content)
		if file.Size > 1_000_000 { // > 1MB
			score += 10
		}

		scored = append(scored, scoredFile{name: file.Name, score: score})
	}

	// SAFETY NET: If no PDF/Docx found at all, try to return ANY large file
	// This prevents the "Mixed Domain" hallucination if the folder has no PDFs
	if len(scored) == 0 && len(files) > 0 {
		// Just grab the largest files regardless of extension
		for _, file := range files {
			if file.Size > 500_000 { // > 500KB
				scored = append(scored, scoredFile{name: file.Name, score: 0})
			}
		}
	}

	// Sort by Score Descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Extract top K results
	var result []string
	for i := 0; i < len(scored) && i < limit; i++ {
		result = append(result, scored[i].name)
	}

	return result
}

func detectTechStack(files []FileSummary) []string {
	stacks := make(map[string]bool)

	for _, file := range files {
		name := strings.ToLower(file.Name)
		ext := strings.ToLower(file.Extension)

		if name == "pubspec.yaml" {
			stacks["Flutter"] = true
		}
		if name == "package.json" {
			stacks["Node.js/JavaScript"] = true
		}
		if name == "go.mod" {
			stacks["Go"] = true
		}
		if name == "requirements.txt" || name == "pipfile" {
			stacks["Python"] = true
		}
		if name == "cargo.toml" {
			stacks["Rust"] = true
		}
		if name == "pom.xml" || name == "build.gradle" {
			stacks["Java"] = true
		}
		if name == "cmakelists.txt" || name == "makefile" {
			stacks["C/C++ (Make/CMake)"] = true
		}
		if ext == ".cpp" || ext == ".c" || ext == ".h" || ext == ".hpp" {
			stacks["C/C++"] = true
		}
		if ext == ".cu" {
			stacks["CUDA (NVIDIA)"] = true
		}
	}

	result := []string{}
	for stack := range stacks {
		result = append(result, stack)
	}
	return result
}

// GeneratePrompt creates the prompt based on your new "Smart Summary" format
// using the ENRICHED DirectorySummary to access file metadata.
func GeneratePrompt(insight *ContentInsight, summary *DirectorySummary) string {

	systemPrompt := `You are Scout, an intelligent directory analyst.

### INSTRUCTIONS:
1. **Analyze** the "stats" and "total_files" for the "This folder contains" section.
2. **Infer** the "Likely Purpose".
3. **Select** interesting "Highlights".
4. **Suggest** actionable next steps.

### ⛔ NEGATIVE CONSTRAINTS (CRITICAL):
- DO NOT output "Step 1", "Step 2", or "Here is the analysis".
- DO NOT output internal reasoning or chain-of-thought.
- OUTPUT ONLY the final result starting with the 📁 emoji.
- DO NOT use Markdown headers (like ## or ###). Use the emojis as headers.

### REQUIRED OUTPUT FORMAT:

📁 This folder contains:
  - [total_files] files total

🎯 Likely Purpose:
  [Hypothesis based strictly on file previews]

🔍 Highlights:
  - [Content Insight]
  - [Common Technology Used]
  - [Key Pattern]

⚠️ Suggestions:
  - [Actionable advice]
  - [Reading recommendation]

### RULES:
- BE TRUTHFUL.
- Keep it concise.`

	// 2. Prepare the Context (The Ingredients)
	type KeyFileContext struct {
		Name     string         `json:"name"`
		Type     string         `json:"extension"`
		Size     string         `json:"size_formatted"`
		Metadata map[string]any `json:"metadata"` // Contains "preview" (code snippets)
	}

	var keyFilesCtx []KeyFileContext

	// Match insight.KeyFiles (names) to summary.Files (data) to get the Metadata
	for _, filename := range insight.KeyFiles {
		for _, f := range summary.Files {
			if f.Name == filename {
				keyFilesCtx = append(keyFilesCtx, KeyFileContext{
					Name:     f.Name,
					Type:     f.Extension,
					Size:     formatBytes(f.Size), // Uses helper function below
					Metadata: f.Metadata,          // <--- Passes code snippets to LLM
				})
				break
			}
		}
	}

	// 3. Build the JSON Payload
	contextData := map[string]any{
		"stats":           insight.FilesByCategory,
		"total_files":     summary.FileCount, // <--- Explicit Total Count to fix Math issue
		"domain_detected": insight.Domain,
		"key_files_data":  keyFilesCtx, // Only the top relevant files with content
		"topics":          insight.Topics,
	}

	contextJSON, _ := json.MarshalIndent(contextData, "", "  ")

	userPrompt := fmt.Sprintf("Analyze this Directory Data:\n%s", string(contextJSON))

	// Llama 3 Prompt Format
	return fmt.Sprintf("<|begin_of_text|><|start_header_id|>system<|end_header_id|>\n\n%s<|eot_id|><|start_header_id|>user<|end_header_id|>\n\n%s<|eot_id|><|start_header_id|>assistant<|end_header_id|>\n\n",
		systemPrompt, userPrompt)
}

// Helper to format bytes (Add this at the bottom of analysis.go if not in utils)
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
