package input

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Knoxianes/go-vim/internal/app"
	"github.com/Knoxianes/go-vim/internal/buffer"
	"github.com/Knoxianes/go-vim/internal/command"
)

var InputChannel chan byte

func LoadInput() {
	reader := bufio.NewReader(os.Stdin)
	for app.AppRunning {
		char, err := reader.ReadByte()
		if err != nil {
			panic("Error reading byte")
		}
		InputChannel <- char
	}
}

func HandleInput() {
	app.CurrentBuffer.PrintBuffer()
	InputChannel = make(chan byte)
	go LoadInput()
	for app.AppRunning {
		switch app.CurrentBuffer.Mode {
		case buffer.NormalMode:
			HandleInputNormalMode()
		case buffer.InsertMode:
			HandleInputInsertMode()
		}
		app.CurrentBuffer.PrintBuffer()
	}
}
func HandleInputInsertMode() {
	select {
	case char, ok := <-InputChannel:
		if !ok {
			app.AppRunning = false
			return
		}
		if char > 31 && char < 127 {
			app.CurrentBuffer.InsertChar(char)
		}
		if char == 27 {
			app.CurrentBuffer.Mode = buffer.NormalMode
		}
		if char == 127 {
			app.CurrentBuffer.DeleteChar()
		}
		if char == 13 {
			app.CurrentBuffer.BreakLine()
		}
	}

}
func HandleInputNormalMode() {
	select {
	case char, ok := <-InputChannel:
		if !ok {
			app.AppRunning = false
			return
		}
		command.ExecuteCommand(fmt.Sprint(int(char)))
	}
}
