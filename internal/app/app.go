package app

import (
	"os/exec"

	"bufio"
	"github.com/Knoxianes/go-vim/internal/buffer"
	"github.com/Knoxianes/go-vim/pkg/terminal"
	"os"
)

var copyRegister string
var opennedBuffers map[string]*buffer.Buffer
var CurrentBuffer *buffer.Buffer
var AppRunning bool

func GetCopyRegister() string {
	return copyRegister
}
func SetCopyRegister(r string) {
	copyRegister = r
}
func NewBuffer(path string) {
	if val, ok := opennedBuffers[path]; ok {
		CurrentBuffer = val
		return
	}
	tmpBuffer := buffer.NewBuffer(path)
	if tmpBuffer == nil {
		return
	}
	opennedBuffers[path] = tmpBuffer
	CurrentBuffer = tmpBuffer
}
func AppInit() {
	setTerminalToRaw()
	defer cleanUp()

	copyRegister = ""
	opennedBuffers = make(map[string]*buffer.Buffer)
	AppRunning = true
	terminal.TerminalSize = GetTerminalSize()

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
	for AppRunning {
		reader := bufio.NewReader(os.Stdin)
		char, err := reader.ReadByte()
		if err != nil {
			panic("Error reading byte")
		}
		if currBuffer.Mode == buffer.NormalMode {
			switch char {
			case 'q':
				AppRunning = false
			case 'j':
				currBuffer.MoveCursorDown()
			case 'k':
				currBuffer.MoveCursorUp()
			case 'h':
				currBuffer.MoveCursorLeft()
			case 'l':
				currBuffer.MoveCursorRight()
			case 'i':
				currBuffer.Mode = buffer.InsertMode
			case 's':
				currBuffer.SaveBuffer()
			}
		} else {
			if char > 31 && char < 127 {
				currBuffer.InsertChar(char)
			}
			if char == 27 {
				currBuffer.Mode = buffer.NormalMode
			}
			if char == 127 {
				currBuffer.DeleteChar()
			}
			if char == 13 {
				currBuffer.BreakLine()
			}
		}

		currBuffer.PrintBuffer()
		print(char)
	}
}

func GetTerminalSize() terminal.Size {
	col, row, err := terminal.GetTerminalSize()
	if err != nil {
		panic("Error getting terminal size")
	}
	return terminal.Size{Col: col, Row: row}
}

func cleanUp() {
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	exec.Command("stty", "-F", "/dev/tty", "-raw").Run()
	terminal.ShowCursor()
	terminal.ClearScreen()
}
func setTerminalToRaw() {
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	exec.Command("stty", "-F", "/dev/tty", "raw").Run()
	terminal.HideCursor()
}
