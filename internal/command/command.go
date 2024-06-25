package command

import (
	"strconv"
	"strings"

	"github.com/Knoxianes/go-vim/internal/app"
	"github.com/Knoxianes/go-vim/internal/buffer"
)

var CommandMap map[string]func()
var Leader = 32

func LoadCommands() {
	CommandMap = make(map[string]func())

	CommandMap["119"] = func() {
		app.CurrentBuffer.SaveBuffer()
	}

	CommandMap["113"] = func() {
		app.AppRunning = false
	}

	CommandMap["106"] = app.CurrentBuffer.MoveCursorDown
	CommandMap["107"] = app.CurrentBuffer.MoveCursorUp
	CommandMap["104"] = app.CurrentBuffer.MoveCursorLeft
	CommandMap["108"] = app.CurrentBuffer.MoveCursorRight

	CommandMap["105"] = func() {
		app.CurrentBuffer.Mode = buffer.InsertMode
	}

}

func ExecuteCommand(command string) {

	cmd, ok := CommandMap[command]
	if ok {
		cmd()
	}
}

func ConvertCommandToString(command string) string {
	var ret string
	splitedCommand := strings.Split(command, " ")
	for _, subString := range splitedCommand {
		asciiNumber, err := strconv.Atoi(subString)
		if err != nil {
			return ""
		}
		if asciiNumber == Leader {
			ret += "leader"
		} else {
			ret += string(byte(asciiNumber))
		}
	}
	return ret
}
