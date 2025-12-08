#!/bin/bash

# Get the directory where this script is located
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export YZMA_LIB="$DIR/.scout/llama"

if [[ "$OSTYPE" == "darwin"* ]]; then
    export DYLD_LIBRARY_PATH="$DIR/.scout/llama"
else
    export LD_LIBRARY_PATH="$DIR/.scout/llama"
fi

"$DIR/bin/scout" "$@"