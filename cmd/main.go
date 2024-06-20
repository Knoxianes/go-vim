package main

import (
	"os"
	"os/exec"

	"github.com/Knoxianes/go-vim/internal/buffer"
)

func main() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	args := os.Args
	if len(args) != 2 {
		panic("Invalid number of arguments")
	}
	var currBuffer *buffer.Buffer
	if args[1][0] == '/' {
		currBuffer = buffer.NewBuffer(args[1])
	} else {
		currentDir, err := os.Getwd()
		if err != nil {
			panic("Error getting current directory")
		}
		currBuffer = buffer.NewBuffer(currentDir + "/" + args[1])
	}
	if currBuffer == nil {
		panic("Error creating buffer")
	}

	currBuffer.PrintBuffer()
	for {
	}
}
