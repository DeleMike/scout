package shell

import (
	"os"
	"os/exec"
)

func (shell *Shell) runExternal(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}
