package main

import (
	"bufio"
	"os"
	"os/exec"

	"github.com/Knoxianes/go-vim/internal/buffer"
	"github.com/Knoxianes/go-vim/pkg/terminal"
)

func main() {

	initApp()

	defer cleanUp()

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
		reader := bufio.NewReader(os.Stdin)
		char, err := reader.ReadByte()
		if err != nil {
			panic("Error reading byte")
		}
		if currBuffer.GetMode() == buffer.NormalMode {
			switch char {
			case 'q':
				return
			case 'j':
				currBuffer.MoveCursorDown()
			case 'k':
				currBuffer.MoveCursorUp()
			case 'h':
				currBuffer.MoveCursorLeft()
			case 'l':
				currBuffer.MoveCursorRight()
			case 'i':
				currBuffer.SetInsertMode()
			case 's':
				currBuffer.SaveBuffer()
			}
		} else {
			if char > 31 && char < 127 {
				currBuffer.InsertChar(char)
			}
			if char == 27 {
				currBuffer.SetNormalMode()
			}
			if char == 127 {
				currBuffer.DeleteChar()
			}
			if char == 13 {
				currBuffer.BreakLine()
			}
		}

		currBuffer.PrintBuffer()
		println(char)
	}
}

func initApp() {

	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	exec.Command("stty", "-F", "/dev/tty", "raw").Run()
	terminal.HideCursor()
}

func cleanUp() {
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	exec.Command("stty", "-F", "/dev/tty", "-raw").Run()
	terminal.ShowCursor()
}
