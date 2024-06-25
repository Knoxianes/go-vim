package app

import (
	"os"

	"github.com/Knoxianes/go-vim/internal/buffer"
	"github.com/Knoxianes/go-vim/pkg/terminal"
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

	copyRegister = ""
	opennedBuffers = make(map[string]*buffer.Buffer)
	AppRunning = true
	terminal.TerminalSize = GetTerminalSize()

	args := os.Args
	if len(args) != 2 {
		panic("Invalid number of arguments")
	}
	if args[1][0] == '/' {
		CurrentBuffer = buffer.NewBuffer(args[1])
	} else {
		currentDir, err := os.Getwd()
		if err != nil {
			panic("Error getting current directory")
		}
		CurrentBuffer = buffer.NewBuffer(currentDir + "/" + args[1])
	}
	if CurrentBuffer == nil {
		panic("Error creating buffer")
	}

}

func GetTerminalSize() terminal.Size {
	col, row, err := terminal.GetTerminalSize()
	if err != nil {
		panic("Error getting terminal size")
	}
	return terminal.Size{Col: col, Row: row}
}
