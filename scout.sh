#!/bin/bash

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PROJECT_ROOT="$(dirname "$DIR")"

export YZMA_LIB="$PROJECT_ROOT/.scout/llama"
export SCOUT_MODEL="$PROJECT_ROOT/.scout/model/llama-3.2-3b-instruct-q4_k_m.gguf"

CORE_BIN="$DIR/scout-core"


# echo "🐛 DEBUG: Scout Wrapper"
# echo "  - Core Binary: $CORE_BIN"
# echo "  - Loading Lib: $YZMA_LIB"
# echo "  - Loading Mod: $SCOUT_MODEL"
# ls -lh "$YZMA_LIB"

if [ ! -f "$CORE_BIN" ]; then
    echo "❌ Error: Could not find core binary at: $CORE_BIN"
    exit 1
fi

# 5. Run the core binary, passing all arguments
exec "$CORE_BIN" "$@"