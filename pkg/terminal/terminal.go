package terminal

import (
	"errors"
	"fmt"

	"golang.org/x/term"
)

func GetTerminalSize() (int, int, error) {
	if !term.IsTerminal(0) {
		return 0, 0, errors.New("Not a terminal")
	}
	return term.GetSize(0)
}

func ClearScreen() {

	fmt.Print("\033[H\033[2J")
}

func BackGroundHighlight() {
	fmt.Print("\033[48;5;240m")
}

func CursorColor() {
	fmt.Print("\033[7m")
}
func CursorBlinking() {
	fmt.Print("\033[5m")
}

func ResetScreenAttributes() {
	fmt.Print("\033[0m")
}

func HideCursor() {
	fmt.Print("\033[?25l")
}
func ShowCursor() {
	fmt.Print("\033[?25h")
}
