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
	cursor  cursor
	mode    int
}

type cursor struct {
	Row int
	Col int
}

func NewBuffer(path string) *Buffer {
	dat, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file")
		return nil
	}
	tmpBuffer := &Buffer{
		path:    path,
		content: bytes.Split(dat, []byte("\n")),
		mode:    NormalMode,
		cursor: cursor{
			Row: 0,
			Col: 0,
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
	b.content[b.cursor.Row] = slices.Insert(b.content[b.cursor.Row], b.cursor.Col, c)
	if b.cursor.Col < len(b.content[b.cursor.Row]) {
		b.cursor.Col++
	}
}

func (b *Buffer) InsertNewline() {
	b.content = slices.Insert(b.content, b.cursor.Row+1, []byte{10})
}

func (b *Buffer) BreakLine() {
	b.content = slices.Insert(b.content, b.cursor.Row+1, b.content[b.cursor.Row][b.cursor.Col:])
	b.content[b.cursor.Row] = b.content[b.cursor.Row][:b.cursor.Col]
	b.cursor.Row++
	b.cursor.Col = 0
}

func (b *Buffer) DeleteChar() {
	if b.cursor.Col == 0 {
		if b.cursor.Row == 0 {
			return
		}
		b.cursor.Col = len(b.content[b.cursor.Row-1])
		if b.cursor.Col < 0 {
			b.cursor.Col = 0
		}
		b.content[b.cursor.Row-1] = append(b.content[b.cursor.Row-1], b.content[b.cursor.Row]...)
		b.DeleteLine()
		return
	}
	b.content[b.cursor.Row] = slices.Delete(b.content[b.cursor.Row], b.cursor.Col-1, b.cursor.Col)
	b.cursor.Col--
}

func (b *Buffer) DeleteLine() {
	if len(b.content) == 1 {
		b.content[0] = []byte{}
		return
	}
	b.content = slices.Delete(b.content, b.cursor.Row, b.cursor.Row+1)
	b.cursor.Row--
}

func (b *Buffer) MoveCursorUp() {
	if b.cursor.Row > 0 {
		b.cursor.Row--
	}
	if b.cursor.Col > len(b.content[b.cursor.Row]) {
		b.cursor.Col = len(b.content[b.cursor.Row])
	}
	if len(b.content[b.cursor.Row]) == 0 {
		b.cursor.Col = 0
	}
}
func (b *Buffer) MoveCursorDown() {

	if b.cursor.Row < len(b.content)-1 {
		b.cursor.Row++
	}
	if b.cursor.Col > len(b.content[b.cursor.Row]) {
		b.cursor.Col = len(b.content[b.cursor.Row])
	}
	if len(b.content[b.cursor.Row]) == 0 {
		b.cursor.Col = 0
	}

}
func (b *Buffer) MoveCursorLeft() {
	if b.cursor.Col > 0 {
		b.cursor.Col--
	}
}
func (b *Buffer) MoveCursorRight() {
	if b.cursor.Col < len(b.content[b.cursor.Row]) {
		b.cursor.Col++
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

	for i, line := range b.content {
		for j, c := range line {
			if i == b.cursor.Row && j == b.cursor.Col {
				terminal.CursorColor()
				if b.mode == InsertMode {
					terminal.CursorBlinking()
				}
			}
			fmt.Printf("%c", c)
			terminal.ResetScreenAttributes()
		}
		if b.cursor.Col == len(line) && i == b.cursor.Row && len(line) > 0 {
			terminal.CursorColor()
			if b.mode == InsertMode {
				terminal.CursorBlinking()
			}
			fmt.Printf("%c", 32)
			terminal.ResetScreenAttributes()
		}
		if len(line) == 0 {
			if i == b.cursor.Row {
				terminal.CursorColor()
				if b.mode == InsertMode {
					terminal.CursorBlinking()
				}
			}
			fmt.Printf("%c", 32)
			terminal.ResetScreenAttributes()
		}
		fmt.Printf("\r\n")
	}
	fmt.Println("cursor Row: ", b.cursor.Row, "cursor Col: ", b.cursor.Col)
}

func (b *Buffer) SetVisualMode() {
	b.mode = VisualMode
}
func (b *Buffer) SetInsertMode() {
	b.mode = InsertMode
}
func (b *Buffer) SetNormalMode() {
	b.mode = NormalMode
}
func (b *Buffer) GetMode() int {
	return b.mode
}
