// shell package
package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Shell class design
type Shell struct {
	prompt string
}

// New initializes a Shell application.
func New() *Shell {
	return &Shell{
		prompt: "scout> ",
	}
}

// Start opens and starts shell to receive commands(inputs) and process them to deliver output.
// Shell basically follow a REPL event loop.
func (shell *Shell) Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(shell.prompt)

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		args := strings.Split(line, " ")

		if len(args) == 0 || line == "" {
			continue
		}

		// check for builtin command
		if handled := shell.runBuiltin(args); handled {
			continue
		}

		// check for external commands(e.g git, go, curl, etc)
		shell.runExternal(args)

	}
}
