package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("scout> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		args := strings.Fields(line)

		// implement shell functions
		switch args[0] {
		case "exit":
			fmt.Println("Byeee!")
			return
		case "cd":
			if len(args) < 2 {
				fmt.Println("usage: cd <path>")
				continue
			}
			if err := os.Chdir(args[1]); err != nil {
				fmt.Println("error:", err)
			}

		default:
			fmt.Println("unknown command:", args[0])
		}

	}
}
