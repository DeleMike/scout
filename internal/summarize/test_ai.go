package summarize

// package main

// import (
//     "fmt"
//     "os"
//     "strings"

//     "github.com/hybridgroup/yzma/pkg/llama"
// )

// func main() {
//     // Check environment variables
//     libPath := os.Getenv("YZMA_LIB")
//     if libPath == "" {
//         fmt.Println("Error: YZMA_LIB environment variable not set")
//         os.Exit(1)
//     }

//     modelPath := "/Users/mac/Downloads/llama-3.2-3b-instruct-q4_k_m.gguf"

//     // Check if model exists
//     if _, err := os.Stat(modelPath); os.IsNotExist(err) {
//         fmt.Printf("Error: Model file not found at %s\n", modelPath)
//         os.Exit(1)
//     }

//     llama.Load(libPath)
//     llama.Init()

//     // Suppress llama.cpp logs
//     llama.LogSet(llama.LogSilent())

//     model := llama.ModelLoadFromFile(modelPath, llama.ModelDefaultParams())
//     defer llama.ModelFree(model)

//     lctx := llama.InitFromModel(model, llama.ContextDefaultParams())
//     defer llama.Free(lctx)

//     vocab := llama.ModelGetVocab(model)

//     // Construct prompt in Llama 3.2 format
//     systemPrompt := "You are a helpful AI assistant."
//     userPrompt := "In two sentences, explain the difference between concurrency and parallelism in Go."
//     fullPrompt := fmt.Sprintf("<|begin_of_text|><|start_header_id|>system<|end_header_id|>\n\n%s<|eot_id|><|start_header_id|>user<|end_header_id|>\n\n%s<|eot_id|><|start_header_id|>assistant<|end_header_id|>\n\n",
//         systemPrompt, userPrompt)

//     // Tokenize the prompt
//     tokens := llama.Tokenize(vocab, fullPrompt, true, false)

//     // Create batch
//     llamaTokens := make([]llama.Token, len(tokens))
//     for i, t := range tokens {
//         llamaTokens[i] = llama.Token(t)
//     }
//     batch := llama.BatchGetOne(llamaTokens)

//     // Initialize sampler
//     sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
//     defer llama.SamplerFree(sampler)
//     llama.SamplerChainAdd(sampler, llama.SamplerInitGreedy())

//     // Generate tokens
//     maxTokens := int32(512)
//     var response strings.Builder

//     for pos := int32(0); pos < maxTokens; pos += batch.NTokens {
//         if llama.Decode(lctx, batch) != 0 {
//             break
//         }

//         token := llama.SamplerSample(sampler, lctx, -1)

//         if llama.VocabIsEOG(vocab, token) {
//             break
//         }

//         buf := make([]byte, 128)
//         length := llama.TokenToPiece(vocab, token, buf, 0, false)

//         if length > 0 {
//             response.WriteString(string(buf[:length]))
//         }

//         batch = llama.BatchGetOne([]llama.Token{token})
//     }

//     // Print only the clean answer
//     fmt.Println(strings.TrimSpace(response.String()))
// }
