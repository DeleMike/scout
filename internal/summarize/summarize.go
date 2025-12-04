package summarize

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hybridgroup/yzma/pkg/llama"
)

// Summarize runs Llama inference on a fully-formed prompt and returns the formatted output.
func Summarize(prompt string) (string, error) {
	libPath := os.Getenv("YZMA_LIB")
	if libPath == "" {
		return "", fmt.Errorf("YZMA_LIB environment variable not set")
	}

	modelPath := "/Users/mac/Downloads/llama-3.2-3b-instruct-q4_k_m.gguf"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return "", fmt.Errorf("model file not found at %s", modelPath)
	}

	// Truncate very long prompts for safety
	if len(prompt) > 12000 {
		prompt = prompt[:12000] + "\n...[truncated due to size]..."
	}

	// Initialize Llama
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

	// Tokenize the prompt as-is (no extra system/user wrapping)
	tokens := llama.Tokenize(vocab, prompt, false, false)
	batchLimit := int(ctxParams.NBatch)
	if len(tokens) > batchLimit {
		tokens = tokens[:batchLimit-1]
	}

	llamaTokens := make([]llama.Token, len(tokens))
	for i, t := range tokens {
		llamaTokens[i] = llama.Token(t)
	}

	batch := llama.BatchGetOne(llamaTokens)

	// Sampler
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

// FormatForTerminal adds simple terminal colors & bold for readability.
func FormatForTerminal(text string) string {
	const (
		Reset = "\033[0m"
		Bold  = "\033[1m"
		Cyan  = "\033[36m"
	)

	// Headers (###) → Cyan & Bold
	headerRe := regexp.MustCompile(`(?m)^###\s*(.*)$`)
	text = headerRe.ReplaceAllString(text, Cyan+Bold+"$1"+Reset)

	// **Bold** → terminal bold
	boldRe := regexp.MustCompile(`\*\*(.*?)\*\*`)
	text = boldRe.ReplaceAllString(text, Bold+"$1"+Reset)

	return text
}
