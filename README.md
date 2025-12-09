# Scout

> A Universal Directory Onboarding. Instantly understand any folder.

[![Go Report Card](https://goreportcard.com/badge/github.com/DeleMike/scout?style=for-the-badge)](https://goreportcard.com/report/github.com/DeleMike/scout)
[![Go Version](https://img.shields.io/github/go-mod/go-version/DeleMike/scout?style=for-the-badge&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](./LICENSE)

![Project Logo](mascot.png)

---

## 📖 About

**Scout** is an AI-powered directory analyzer that turns any folder into a readable story.  
It scans files, detects the domain (Software, Medical, Legal, Creative, etc.), and generates a clean summary using a **local Llama 3.2 model** — ensuring zero cloud dependency.

Use Scout when:
- inheriting a legacy project  
- joining a new codebase  
- cleaning up messy folders  
- onboarding teammates  
- reviewing documents or reports  

Scout extracts structure, intent, and “where to start” guidance so you don’t have to dig manually.

---

## ✨ Features

- **🧠 Domain-aware insights:** Recognizes codebases, financial docs, creative assets, research folders, and more.
- **🔒 Local-first & private:** Runs on your machine using GGUF models (Llama 3.2). No APIs. No tracking.
- **⚡️ Fast filesystem scanner:** Ignores noise (`node_modules`, `.git`, caches) and extracts useful metadata only.
- **🐚 Built-in interactive shell:** Includes `cd`, `ls`, `pwd`, and the powerful `sc` analyzer command.
- **📄 Multi-format extraction:** Reads previews from PDFs, DOCX, Markdown, spreadsheets, images, and code.

---

## 🚀 Getting Started

Scout includes a **Makefile** that automates setup, verification, building, and running.  
This is the recommended way to install and use Scout locally. For a manual alternative, see the [Detailed Installation](#-detailed-installation-manual-alternative) section below.

### Prerequisites
- **Go 1.21+** (install via [golang.org/dl](https://golang.org/dl) or your package manager, e.g., `brew install go`).
- **A C compiler** (GCC or Clang for CGO):  
  - macOS: `xcode-select --install` or `brew install gcc`.  
  - Linux (Ubuntu/Debian): `sudo apt install build-essential`.  
  - Windows: Use WSL
- **Llama 3.2 GGUF model** (recommended: `llama-3.2-3b-instruct-q4_k_m.gguf`).  
  Download from: [Hugging Face](https://huggingface.co/hugging-quants/Llama-3.2-3B-Instruct-Q4_K_M-GGUF).

### 1. Clone the Repository
```bash
git clone https://github.com/DeleMike/scout.git
cd scout
```

### 2. First-Time Setup
```bash
make setup
```
This creates the required directory structure:  
```
.scout/
  llama/   # place libllama.dylib (macOS) or libllama.so (Linux) here
  model/   # place your GGUF model here
```

### 3. Install Runtime Libraries
Scout relies on [llama.cpp](https://github.com/ggerganov/llama.cpp) for local inference. Install the shared library (`libllama.dylib` on macOS or `libllama.so` on Linux):  

**Go Integration (via YZMA):**  
```bash
go get github.com/hybridgroup/yzma/pkg/llama@latest
go mod tidy
```
*(YZMA bundles libs; follow its docs to extract `libllama` if needed.)*

Copy the library to `.scout/llama/`.  
Copy the downloaded GGUF model to `.scout/model/`.

### 4. Verify Installation
```bash
make check
```
This confirms:  
- ✔️ Model exists  
- ✔️ Llama runtime library exists  

If issues arise (e.g., "Library not loaded"), set library paths:  
- macOS: `export DYLD_LIBRARY_PATH="$(pwd)/.scout/llama:$DYLD_LIBRARY_PATH"`  
- Linux: `export LD_LIBRARY_PATH="$(pwd)/.scout/llama:$LD_LIBRARY_PATH"`

### 5. Build Scout
```bash
make build
```
Outputs:  
- `bin/scout-core` (core binary)  
- `bin/scout` (wrapper script)

### 6. Run Scout
```bash
make run
```
Or directly:  
```bash
./bin/scout
```

#### Global Installation (Optional)
```bash
sudo mv bin/scout /usr/local/bin/
```
Now run `scout` from anywhere.

#### Set Environment Variables (for Persistence)
Add to `~/.zshrc` or `~/.bashrc`:  
```bash
export YZMA_LIB="$(pwd)/.scout/llama"
export DYLD_LIBRARY_PATH="$(pwd)/.scout/llama"  # macOS
# export LD_LIBRARY_PATH="$(pwd)/.scout/llama"  # Linux
export SCOUT_MODEL="$(pwd)/.scout/model/llama-3.2-3b-instruct-q4_k_m.gguf"
```
Verify: `echo $YZMA_LIB && echo $SCOUT_MODEL`.

---

## 🔧 Detailed Installation (Manual Alternative)
If you prefer not to use the Makefile:  

### Install Runtime Libraries
Follow the steps in [Step 3](#3-install-runtime-libraries) above.

### Install Scout
```bash
go mod tidy
CGO_ENABLED=1 go build -o scout ./cmd/scout/main.go
```

### Prepare Model & Runtime
Manually create:  
```bash
mkdir -p .scout/model .scout/llama
```
Place files as described in [Step 2](#2-first-time-setup).

### Set Environment Variables
As in [Step 6](#6-run-scout) above.

---

## 🛠 Usage

### Interactive Shell
Launch Scout to enter the shell:  
```bash
./bin/scout
```
Commands:  
```bash
scout> ls                 # List files
scout> cd ../legacy-code  # Navigate directories
scout> scout                 # Analyze the current folder
scout> sc ./frontend      # Analyze a specific subfolder
```

### Quick Scan (Headless Mode)
Run Scout directly from your terminal to scan a folder and exit immediately. Perfect for quick checks.
```bash
# Scan current directory
scout .
sc "."

# Scan a specific path
sc /Users/dev/projects/My-Go-Project
scout "/Users/dev/projects/My-Go-Project"
```

### Saving Reports (Export to File)
Need to share the analysis? Pipe the output to a text file
```bash
scout . >> analysis_report.txt

sc "/Users/dev/projects/My-Go-Project" >> report.txt
```

---

## 🧪 Example Output

```text
🔎 Scouting: /Users/mac/projects/My-Go-Project
✅ Found 42 files (88% confidence: Software Domain)
🤖 Generating AI insights...

================================================================================
📁 This folder contains:
  - 18 Go source files
  - 5 YAML configuration files
  - 42 files total

🎯 Likely Purpose:
  A high-throughput logging service written in Go, designed to ingest and 
  store event data efficiently.

🔍 Highlights:
  - Common Technology: Go (Backend), SQLite (Storage), Docker (Deployment)
  - Key Pattern: Uses a worker-pool pattern for concurrent log processing.
  - Content Insight: "main.go" initializes a TCP server on port 9000.

👀 Suggestions:
  - Start by reading "internal/ingest/worker.go" to understand the concurrency model.
  - Check "docker-compose.yml" to see how the database is orchestrated.
================================================================================
```

---

## 🤝 Contributing

We welcome contributions, especially around:  
- new extractors (PPTX, EPUB, media metadata)  
- better domain heuristics  
- performance improvements  
- Windows support for local LLMs  

Steps:  
1. Fork the repo  
2. `git checkout -b feature/my-feature`  
3. Commit + push  
4. Open a PR  

---

## 📜 License

Licensed under the MIT License. See the [LICENSE](./LICENSE) file.