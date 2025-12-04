package summarize

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func Summarize(prompt string) (string, error) {
	libPath := os.Getenv("YZMA_LIB")
	if libPath == "" {
		return "", fmt.Errorf("YZMA_LIB environment variable not set")
	}

	modelPath := "/Users/mac/Downloads/llama-3.2-3b-instruct-q4_k_m.gguf"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return "", fmt.Errorf("model file not found at %s", modelPath)
	}

	llama.Load(libPath)
	llama.Init()
	llama.LogSet(llama.LogSilent())

	model := llama.ModelLoadFromFile(modelPath, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	ctxParams := llama.ContextDefaultParams()
	// 1. INCREASE CONTEXT WINDOW
	// 22k files need space. 16k tokens is safe for Mac M1/M2/M3.
	ctxParams.NCtx = 16384
	ctxParams.NBatch = 4096

	lctx := llama.InitFromModel(model, ctxParams)
	defer llama.Free(lctx)

	vocab := llama.ModelGetVocab(model)

	// 2. TOKENIZE FULL PROMPT (No Truncation!)
	tokens := llama.Tokenize(vocab, prompt, false, false)

	// fmt.Printf("📊 Token Count: %d / %d\n", len(tokens), ctxParams.NCtx)

	if len(tokens) > int(ctxParams.NCtx) {
		return "", fmt.Errorf("prompt is too large (%d tokens). Limit is %d", len(tokens), ctxParams.NCtx)
	}

	// 3. CHUNKED DECODING (Prefill)
	// This prevents the "GGML_ASSERT" crash and handles long prompts gracefully.
	batchSize := int(ctxParams.NBatch)

	for i := 0; i < len(tokens); i += batchSize {
		end := i + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}

		chunk := tokens[i:end]

		chunkLlama := make([]llama.Token, len(chunk))
		for j, t := range chunk {
			chunkLlama[j] = llama.Token(t)
		}

		batch := llama.BatchGetOne(chunkLlama)
		if llama.Decode(lctx, batch) != 0 {
			return "", fmt.Errorf("llama decode failed on prompt chunk %d-%d", i, end)
		}
	}

	// 4. GENERATION LOOP
	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	defer llama.SamplerFree(sampler)
	llama.SamplerChainAdd(sampler, llama.SamplerInitGreedy())

	maxTokens := int32(1024)
	var response strings.Builder

	// Prime the pump with a dummy sample to get started
	token := llama.SamplerSample(sampler, lctx, -1)

	// We start generating *after* the full prompt has been consumed
	batch := llama.BatchGetOne([]llama.Token{token})

	// Check if the very first token is useful
	buf := make([]byte, 128)
	length := llama.TokenToPiece(vocab, token, buf, 0, false)
	if length > 0 {
		response.WriteString(string(buf[:length]))
	}

	for pos := int32(0); pos < maxTokens; pos++ {
		if llama.Decode(lctx, batch) != 0 {
			break
		}

		token = llama.SamplerSample(sampler, lctx, -1)
		if llama.VocabIsEOG(vocab, token) {
			break
		}

		length = llama.TokenToPiece(vocab, token, buf, 0, false)
		if length > 0 {
			piece := string(buf[:length])
			if strings.Contains(piece, "<|") || strings.Contains(piece, "assistant<|") {
				break
			}
			response.WriteString(piece)
		}

		batch = llama.BatchGetOne([]llama.Token{token})
	}

	rawOutput := strings.TrimSpace(response.String())
	cleanOutput := FormatForTerminal(rawOutput)

	return cleanOutput, nil
}

// FormatForTerminal adds colors based on the Emojis used in the Scout output
func FormatForTerminal(text string) string {
	const (
		Reset  = "\033[0m"
		Bold   = "\033[1m"
		Cyan   = "\033[36m"
		Green  = "\033[32m"
		Yellow = "\033[33m"
		Red    = "\033[31m"
	)

	lines := strings.Split(text, "\n")
	var formatted []string

	for _, line := range lines {
		// Colorize based on Emojis
		if strings.Contains(line, "📁") {
			line = Cyan + Bold + line + Reset
		} else if strings.Contains(line, "🎯") {
			line = Green + Bold + line + Reset
		} else if strings.Contains(line, "🔍") {
			line = Yellow + Bold + line + Reset
		} else if strings.Contains(line, "⚠️") || strings.Contains(line, "👀") {
			line = Red + Bold + line + Reset
		} else if strings.Contains(line, "Step") && strings.Contains(line, ":") {
			line = Bold + line + Reset
		}

		// Handle bolding **text**
		boldRe := regexp.MustCompile(`\*\*(.*?)\*\*`)
		line = boldRe.ReplaceAllString(line, Bold+"$1"+Reset)

		formatted = append(formatted, line)
	}

	return strings.Join(formatted, "\n")
}
