package buffer

import (
	"bytes"
	"fmt"
	"os"
	"slices"

	"github.com/Knoxianes/go-vim/pkg/terminal"
)

const NormalMode = 0
const InsertMode = 1
const VisualMode = 2

type Buffer struct {
	path    string
	content [][]byte
	Cursor  Cursor
	Mode    int
	View    View
}

type View struct {
	LowRow  int
	HighRow int
	LowCol  int
	HighCol int
}

type Cursor struct {
	Row int
	Col int
}

func NewBuffer(path string) *Buffer {

	dat, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Buffer{
				path:    path,
				content: [][]byte{{}},
				Mode:    NormalMode,
				Cursor: Cursor{
					Row: 0,
					Col: 0,
				},
				View: View{
					LowRow:  0,
					HighRow: terminal.TerminalSize.Row - 2,
					LowCol:  0,
					HighCol: terminal.TerminalSize.Col,
				},
			}
		}
		fmt.Println("Error reading file")
		return nil
	}
	tmpBuffer := &Buffer{
		path:    path,
		content: bytes.Split(dat, []byte("\n")),
		Mode:    NormalMode,
		Cursor: Cursor{
			Row: 0,
			Col: 0,
		},
		View: View{
			LowRow:  0,
			HighRow: terminal.TerminalSize.Row - 2,
			LowCol:  0,
			HighCol: terminal.TerminalSize.Col,
		},
	}
	tmpBuffer.ConvertTabsToSpaces()
	return tmpBuffer
}

func (b *Buffer) ConvertTabsToSpaces() {
	for i := range b.content {
		b.content[i] = bytes.ReplaceAll(b.content[i], []byte{9}, []byte{32, 32, 32, 32})
	}

}

func (b *Buffer) ConvertSpacesToTabs() {
	for i := range b.content {
		b.content[i] = bytes.ReplaceAll(b.content[i], []byte{32, 32, 32, 32}, []byte{9})
	}

}

func (b *Buffer) InsertChar(c byte) {
	b.content[b.Cursor.Row] = slices.Insert(b.content[b.Cursor.Row], b.Cursor.Col, c)
	if b.Cursor.Col < len(b.content[b.Cursor.Row]) {
		b.Cursor.Col++
	}
	if b.Cursor.Col > b.View.HighCol {
		b.View.LowCol++
		b.View.HighCol++
	}
}

func (b *Buffer) InsertNewline() {
	b.content = slices.Insert(b.content, b.Cursor.Row+1, []byte{10})
}

func (b *Buffer) BreakLine() {
	b.content = slices.Insert(b.content, b.Cursor.Row+1, b.content[b.Cursor.Row][b.Cursor.Col:])
	b.content[b.Cursor.Row] = b.content[b.Cursor.Row][:b.Cursor.Col]
	b.Cursor.Row++
	b.Cursor.Col = 0
	if b.Cursor.Row > b.View.HighRow {
		b.View.LowRow++
		b.View.HighRow++
	}
}

func (b *Buffer) DeleteChar() {
	if b.Cursor.Col == 0 {
		if b.Cursor.Row == 0 {
			return
		}
		b.Cursor.Col = len(b.content[b.Cursor.Row-1])
		if b.Cursor.Col < 0 {
			b.Cursor.Col = 0
		}
		b.content[b.Cursor.Row-1] = append(b.content[b.Cursor.Row-1], b.content[b.Cursor.Row]...)
		b.DeleteLine()
		return
	}
	b.content[b.Cursor.Row] = slices.Delete(b.content[b.Cursor.Row], b.Cursor.Col-1, b.Cursor.Col)
	b.Cursor.Col--
	if b.Cursor.Col < b.View.LowCol {
		b.View.LowCol--
		b.View.HighCol--
	}
}

func (b *Buffer) DeleteLine() {
	if len(b.content) == 1 {
		b.content[0] = []byte{}
		return
	}
	b.content = slices.Delete(b.content, b.Cursor.Row, b.Cursor.Row+1)
	b.Cursor.Row--
	if b.View.LowRow > 0 {
		b.View.LowRow--
	}
	if b.View.HighRow > terminal.TerminalSize.Row-2 {
		b.View.HighRow--
	}
}

func (b *Buffer) MoveCursorUp() {
	if b.Cursor.Row > 0 {
		b.Cursor.Row--
	}
	if b.Cursor.Col > len(b.content[b.Cursor.Row]) {
		b.Cursor.Col = len(b.content[b.Cursor.Row])
	}
	if len(b.content[b.Cursor.Row]) == 0 {
		b.Cursor.Col = 0
	}
	if b.Cursor.Row < b.View.LowRow {
		b.View.LowRow--
		b.View.HighRow--
	}
}
func (b *Buffer) MoveCursorDown() {

	if b.Cursor.Row < len(b.content)-1 {
		b.Cursor.Row++
	}
	if b.Cursor.Col > len(b.content[b.Cursor.Row]) {
		b.Cursor.Col = len(b.content[b.Cursor.Row])
	}
	if len(b.content[b.Cursor.Row]) == 0 {
		b.Cursor.Col = 0
	}
	if b.Cursor.Row > b.View.HighRow {
		b.View.LowRow++
		b.View.HighRow++
	}

}
func (b *Buffer) MoveCursorLeft() {
	if b.Cursor.Col > 0 {
		b.Cursor.Col--
	}
	if b.Cursor.Col < b.View.LowCol {
		b.View.LowCol--
		b.View.HighCol--
	}
}
func (b *Buffer) MoveCursorRight() {
	if b.Cursor.Col < len(b.content[b.Cursor.Row]) {
		b.Cursor.Col++
	}
	if b.Cursor.Col > b.View.HighCol {
		b.View.LowCol++
		b.View.HighCol++
	}
}
func (b *Buffer) SaveBuffer() {
	file, err := os.Create(b.path)
	if err != nil {
		fmt.Println("Error saving file")
		return
	}

	defer file.Close()

	b.ConvertSpacesToTabs()
	for _, line := range b.content {
		file.Write(append(line, 10))
	}
	b.ConvertTabsToSpaces()

}
func (b *Buffer) PrintBuffer() {
	terminal.ClearScreen()
	var numberOfLinesPrinted = 0
	for i, line := range b.content {
		if i < b.View.LowRow || i > b.View.HighRow {
			continue
		}

		PrintLineNumber(i+1, len(b.content))
		for j, c := range line {
			if j < b.View.LowCol && j > b.View.HighCol {
				continue
			}
			if i == b.Cursor.Row && j == b.Cursor.Col {
				terminal.CursorColor()
				if b.Mode == InsertMode {
					terminal.CursorBlinking()
				}
			}
			fmt.Printf("%c", c)
			terminal.ResetScreenAttributes()
		}
		if b.Cursor.Col == len(line) && i == b.Cursor.Row && len(line) > 0 {
			terminal.CursorColor()
			if b.Mode == InsertMode {
				terminal.CursorBlinking()
			}
			fmt.Printf("%c", 32)
			terminal.ResetScreenAttributes()
		}
		if len(line) == 0 {
			if i == b.Cursor.Row {
				terminal.CursorColor()
				if b.Mode == InsertMode {
					terminal.CursorBlinking()
				}
			}
			fmt.Printf("%c", 32)
			terminal.ResetScreenAttributes()
		}
		numberOfLinesPrinted++
		fmt.Printf("\r\n")
	}
	if numberOfLinesPrinted <= b.View.HighRow-b.View.LowRow {
		for i := numberOfLinesPrinted; i <= b.View.HighRow-b.View.LowRow; i++ {
			fmt.Printf("\r\n")
		}
	}
	fmt.Print("Cursor Row: ", b.Cursor.Row, " Cursor Col: ", b.Cursor.Col, " LowRow: ", b.View.LowRow, " HighRow: ", b.View.HighRow, " LowCol: ", b.View.LowCol, " HighCol: ", b.View.HighCol, " Mode: ", b.Mode)
}

func PrintLineNumber(lineNumber int, maxLineNumber int) {
	numberOfDigits := FindNumberOfDigits(maxLineNumber)
	for i := 0; i < numberOfDigits-FindNumberOfDigits(lineNumber); i++ {
		fmt.Print(" ")
	}
	fmt.Print(lineNumber)
	fmt.Print(" ")
}
func FindNumberOfDigits(number int) int {
	ret := 0
	for number > 0 {
		number = number / 10
		ret++
	}
	return ret
}
