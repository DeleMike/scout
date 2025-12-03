package summarize

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func Summarize(jsonSummary string) (string, error) {
	libPath := os.Getenv("YZMA_LIB")
	if libPath == "" {
		return "", fmt.Errorf("YZMA_LIB environment variable not set")
	}

	modelPath := "/Users/mac/Downloads/llama-3.2-3b-instruct-q4_k_m.gguf"

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return "", fmt.Errorf("model file not found at %s", modelPath)
	}

	if len(jsonSummary) > 12000 {
		jsonSummary = jsonSummary[:12000] + "\n...[truncated due to size]..."
	}

	llama.Load(libPath)
	llama.Init()

	llama.LogSet(llama.LogSilent())

	model := llama.ModelLoadFromFile(modelPath, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 4096
	ctxParams.NBatch = 4096
	lctx := llama.InitFromModel(model, ctxParams)

	defer llama.Free(lctx)

	vocab := llama.ModelGetVocab(model)

	systemPrompt := `You are Scout, an intelligent directory analyst.
Your job is to orient a human user who has just opened this folder.

Analyze the file names, extensions, and folder structure in the JSON to determine the "Domain" of this folder.

### 1. DETECT THE DOMAIN & ADAPT YOUR TONE:
- **Software Project:** Act like a Senior Engineer. Identify the stack (e.g., Flutter, Go, Python), architecture (MVC, Clean), and key modules (Auth, Payments, UI).
- **Documents/Admin:** Act like an Executive Assistant. Summarize the topics (Taxes, Medical, Legal) and urgency.
- **Media/Photos:** Act like a friendly Librarian. Be casual. "Just enjoy the memories!"
- **Messy/Downloads:** Act like a Organizer. Suggest how to clean it up.

### 2. REQUIRED OUTPUT FORMAT:
You must answer these three questions clearly. Use emojis.

**📂 What's going on in here?**
(A concise summary of the contents and the domain.)

**🚀 What can I do with it?**
(The practical utility. If code, explain the features. If docs, explain the knowledge.)

**📍 How do I get started?**
(Bullet points pointing to specific files. For code, find the entry point (main.dart, index.js). For docs, find the most recent or important file.)

**Constraint:** Keep it under 200 words. Be human.`

	userPrompt := fmt.Sprintf("Here is the directory scan:\n%s", jsonSummary)

	fullPrompt := fmt.Sprintf("<|begin_of_text|><|start_header_id|>system<|end_header_id|>\n\n%s<|eot_id|><|start_header_id|>user<|end_header_id|>\n\n%s<|eot_id|><|start_header_id|>assistant<|end_header_id|>\n\n",
		systemPrompt, userPrompt)

	tokens := llama.Tokenize(vocab, fullPrompt, false, false)

	batchLimit := int(ctxParams.NBatch)
	if len(tokens) > batchLimit {
		tokens = tokens[:batchLimit-1]
	}

	llamaTokens := make([]llama.Token, len(tokens))
	for i, t := range tokens {
		llamaTokens[i] = llama.Token(t)
	}

	batch := llama.BatchGetOne(llamaTokens)

	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	defer llama.SamplerFree(sampler)
	llama.SamplerChainAdd(sampler, llama.SamplerInitGreedy())

	maxTokens := int32(1024)
	var response strings.Builder

	for pos := int32(0); pos < maxTokens; pos += batch.NTokens {
		if llama.Decode(lctx, batch) != 0 {
			break
		}

		token := llama.SamplerSample(sampler, lctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			break
		}

		buf := make([]byte, 128)
		length := llama.TokenToPiece(vocab, token, buf, 0, false)

		if length > 0 {
			response.WriteString(string(buf[:length]))
		}

		batch = llama.BatchGetOne([]llama.Token{token})
	}

	rawOutput := strings.TrimSpace(response.String())
	cleanOutput := FormatForTerminal(rawOutput)

	return cleanOutput, nil
}

func FormatForTerminal(text string) string {
	// ANSI Escape Codes
	const (
		Reset = "\033[0m"
		Bold  = "\033[1m"
		Cyan  = "\033[36m" // Nice color for headers
	)

	// 1. Remove "###" headers and make that line Cyan & Bold
	// Regex matches: ### Text -> [Cyan]Text[Reset]
	headerRe := regexp.MustCompile(`(?m)^###\s*(.*)$`)
	text = headerRe.ReplaceAllString(text, Cyan+Bold+"$1"+Reset)

	// 2. Convert **Bold** to actual Terminal Bold
	// Regex matches: **Text** -> [Bold]Text[Reset]
	boldRe := regexp.MustCompile(`\*\*(.*?)\*\*`)
	text = boldRe.ReplaceAllString(text, Bold+"$1"+Reset)

	return text
}
