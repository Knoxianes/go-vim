package main

import (
	"os/exec"

	"github.com/Knoxianes/go-vim/internal/app"
	"github.com/Knoxianes/go-vim/internal/command"
	"github.com/Knoxianes/go-vim/internal/input"
	"github.com/Knoxianes/go-vim/pkg/terminal"
)

func main() {
	setTerminalToRaw()
	defer cleanUp()
	app.AppInit()
	command.LoadCommands()
	input.HandleInput()

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
