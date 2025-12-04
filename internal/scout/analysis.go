package scout

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/DeleMike/scout/internal/scanner"
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
		extractMixedInsights(insight)
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

	if categories["pdf"] > 5 && docPercent > 0.5 {
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

func extractMixedInsights(insight *ContentInsight) {
	insight.Topics = []string{"various files"}
	insight.Recommendations = []string{
		"This looks like a mixed collection",
		"Consider organizing by file type or purpose",
		"Check the largest files first",
	}
}

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
		// only PDFs and DOCX are considered "important docs"
		if file.Extension != ".pdf" && file.Extension != ".docx" {
			continue
		}

		score := 0
		name := strings.ToLower(file.Name)

		// keyword scoring
		if strings.Contains(name, "summary") || strings.Contains(name, "overview") {
			score += 50
		}
		if strings.Contains(name, "final") || strings.Contains(name, "important") {
			score += 40
		}

		// year relevance (2 years ago, last year, this year)
		if strings.Contains(name, years[2]) { // this year
			score += 30
		} else if strings.Contains(name, years[1]) { // last year
			score += 20
		} else if strings.Contains(name, years[0]) { // two years ago
			score += 10
		}

		// size weighting (big files could have sime information for us)
		if file.Size > 1_000_000 { // >1MB
			score += 10
		}

		scored = append(scored, scoredFile{name: file.Name, score: score})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// extract top K(limit) names
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
	}

	result := []string{}
	for stack := range stacks {
		result = append(result, stack)
	}
	return result
}

// GeneratePrompt creates an expressive, fun, domain-aware prompt for the LLM
func GeneratePrompt(insight *ContentInsight, scanResult *scanner.ScanResult) string {
	var systemPrompt string

	switch insight.Domain {
	case DomainStudyMaterials:
		systemPrompt = `You are Scout, a friendly and energetic study buddy 📚✨. 
This folder contains study or exam prep materials.

Your job: Help the user quickly grasp what they have, what subjects/topics, and how to approach them effectively. Be fun, encouraging, and expressive.

FORMAT (use exactly this):
📂 What's here: [One lively sentence summarizing topics and materials]
🚀 What to do: [One sentence on how to tackle or use them effectively]
📍 Start here:
• [filename] - [why start here, in a fun/encouraging tone]
• [filename] - [what's next and why]
• [filename] - [and next, catchy suggestion]

Keep it under 100 words. Use actual filenames, be human, playful, and motivating! 🎯`

	case DomainSoftwareProject:
		systemPrompt = `You are Scout, a clever, friendly senior engineer 🤓💻. 
This folder contains a code project.

Your job: Explain what this project does, its tech stack, key modules, and where to start exploring. Be clear, engaging, and slightly playful, like a good tour guide.

FORMAT (use exactly this):
📂 What's here: [Tech stack, main purpose, key modules in one sentence]
🚀 What to do: [What the app/project does, in a concise and lively way]
📍 Start here:
• [filename or path] - [what it shows / why important, catchy wording]
• [filename or path] - [next key module, fun/expressive guidance]
• [filename or path] - [another important piece]

Keep it under 100 words. Use real filenames. Make it readable and friendly! 🚀`

	case DomainMedia:
		systemPrompt = `You are Scout, a super chill, fun media buddy 😎🎵📸. 
This folder has videos, music, or images.

Your job: Be casual, entertaining, and fun. Encourage browsing, enjoyment, and organization in a friendly tone.

FORMAT (use exactly this):
📂 What's here: [Describe the media, quantity, or type in a fun way]
🚀 What to do: [Encourage enjoying it, maybe organizing or sharing]
📍 Start here:
• Just browse and enjoy! 🎉
• Maybe organize by date/event if you want
• Backup somewhere safe

Keep it light, fun, under 75 words! 😁`

	case DomainDocuments:
		systemPrompt = `You are Scout, a witty, helpful executive assistant 📝✨. 
This folder has documents.

Your job: Explain what's in these docs, how to navigate, and give a few actionable suggestions. Be clear but make it interesting.

FORMAT (use exactly this):
📂 What's here: [Types of documents/topics, catchy summary]
🚀 What to do: [How to use these documents effectively]
📍 Start here:
• [filename] - [why start here, witty or fun phrasing]
• [filename] - [next important doc]
• [organization tip, short and expressive]

Keep it professional but lively, under 100 words.`

	case DomainFinancial:
		systemPrompt = `You are Scout, a friendly financial guide 💰📊. 
This folder contains financial docs.

Your job: Explain the records clearly, suggest order and organization, be concise but lively and helpful.

FORMAT (use exactly this):
📂 What's here: [Types of financial docs]
🚀 What to do: [What info is available and why it matters]
📍 Start here:
• [filename] - [what it contains, fun phrasing optional]
• [organization tip]
• [backup reminder]

Keep it professional but expressive, under 100 words.`

	default:
		systemPrompt = `You are Scout, a friendly organizer for mixed folders 😃📂. 
Your job: Summarize what's in this folder, give actionable next steps, and be lively, catchy, and human.

FORMAT (use exactly this):
📂 What's here: [Brief, expressive summary of contents]
🚀 What to do: [Practical, fun next steps]
📍 Start here:
• [filename or suggestion 1]
• [filename or suggestion 2]
• [filename or suggestion 3]

Keep it under 100 words, clear, and engaging!`

	}

	// Build concise context
	contextData := map[string]any{
		"domain":      string(insight.Domain),
		"confidence":  fmt.Sprintf("%.0f%%", insight.Confidence*100),
		"file_count":  len(scanResult.Files),
		"categories":  insight.FilesByCategory,
		"topics":      insight.Topics,
		"date_range":  insight.DateRange,
		"key_files":   insight.KeyFiles,
		"total_files": len(scanResult.Files),
	}

	// Only include top key files
	maxFiles := 5
	if len(insight.KeyFiles) < maxFiles {
		maxFiles = len(insight.KeyFiles)
	}
	contextData["key_files"] = insight.KeyFiles[:maxFiles]

	contextJSON, _ := json.MarshalIndent(contextData, "", "  ")

	userPrompt := fmt.Sprintf(`Directory data:
%s

Instructions:
- Use the EXACT filenames from key_files
- Include paths if needed
- Follow the format exactly
- Be lively, expressive, fun, and helpful
- Suggest actionable next steps
- For software, highlight entry points and key modules
- For study materials, suggest a study order
- For documents/financial, highlight importance clearly
- Use emojis when it enhances clarity or fun`, string(contextJSON))

	// Use Llama 3.2 prompt format
	return fmt.Sprintf("<|begin_of_text|><|start_header_id|>system<|end_header_id|>\n\n%s<|eot_id|><|start_header_id|>user<|end_header_id|>\n\n%s<|eot_id|><|start_header_id|>assistant<|end_header_id|>\n\n",
		systemPrompt, userPrompt)
}
