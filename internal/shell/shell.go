// shell package
package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Shell represents an interactive command-line interface
// that processes user input in a REPL loop.
type Shell struct {
	prompt string // Command prompt displayed to user
}

// New creates and initializes a new Shell instance
// with the default "scout> " prompt.
//
// Returns:
//   - *Shell: Configured shell ready to accept commands
func New() *Shell {
	return &Shell{
		prompt: "scout> ",
	}
}

// Start begins the REPL (Read-Eval-Print Loop) for the shell.
// It continuously reads user input, parses commands, and executes them
// until the user exits.
//
// The shell supports:
//   - Built-in commands (pwd, ls, sc, exit)
//   - External commands (git, curl, etc.)
//
// The loop continues indefinitely until explicitly terminated.
func (shell *Shell) Start() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(shell.prompt)

		// Read user input until newline
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		// Parse into command and arguments
		args := strings.Split(line, " ")

		if len(args) == 0 || line == "" {
			continue
		}

		// check for builtin command
		if handled := shell.runBuiltin(args); handled {
			continue
		}

		// check for external commands(e.g git, go, curl, etc); can also be a fallback
		shell.runExternal(args)

	}
}
