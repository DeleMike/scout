package shell

import (
	"os"
	"os/exec"
)

// runExternal executes external system commands by spawning
// a new process and connecting it to the shell's I/O streams.
//
// This allows the shell to support any installed system command
// (git, curl, etc.) without implementing them directly.
//
// Parameters:
//   - args: Command and arguments (args[0] is the command name)
//
// The command inherits the shell's stdin, stdout, and stderr,
// allowing interactive commands to work properly.
func (shell *Shell) runExternal(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}
