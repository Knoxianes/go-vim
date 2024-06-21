package buffer

import (
	"bytes"
	"fmt"
	"os"
	"slices"

	"github.com/Knoxianes/go-vim/pkg/terminal"
)

const NormalCursor = 0
const InsertCursor = 1

type Buffer struct {
	Path    string
	Content [][]byte
	Cursor  Cursor
}

type Cursor struct {
	Row  int
	Col  int
	Type int
}

func NewBuffer(path string) *Buffer {
	dat, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file")
		return nil
	}
	tmpBuffer := &Buffer{
		Path:    path,
		Content: bytes.Split(dat, []byte("\n")),
		Cursor: Cursor{
			Row:  0,
			Col:  0,
			Type: NormalCursor,
		},
	}
	tmpBuffer.ConvertTabsToSpaces()
	return tmpBuffer
}

func (b *Buffer) ConvertTabsToSpaces() {
	for i := range b.Content {
		b.Content[i] = bytes.ReplaceAll(b.Content[i], []byte{9}, []byte{32, 32, 32, 32})
	}

}

func (b *Buffer) ConvertSpacesToTabs() {
	for i := range b.Content {
		b.Content[i] = bytes.ReplaceAll(b.Content[i], []byte{32, 32, 32, 32}, []byte{9})
	}

}

func (b *Buffer) InsertChar(c byte) {
	b.Content[b.Cursor.Row] = slices.Insert(b.Content[b.Cursor.Row], b.Cursor.Col, c)
	if b.Cursor.Col < len(b.Content[b.Cursor.Row]) {
		b.Cursor.Col++
	}
}

func (b *Buffer) InsertNewline() {
	b.Content = slices.Insert(b.Content, b.Cursor.Row+1, []byte{10})
}

func (b *Buffer) DeleteChar() {
	if b.Cursor.Col == 0 {
		if b.Cursor.Row == 0 {
			return
		}
		b.Cursor.Col = len(b.Content[b.Cursor.Row-1])
		if b.Cursor.Col < 0 {
			b.Cursor.Col = 0
		}
		b.Content[b.Cursor.Row-1] = append(b.Content[b.Cursor.Row-1], b.Content[b.Cursor.Row]...)
		b.DeleteLine()
		return
	}
	b.Content[b.Cursor.Row] = slices.Delete(b.Content[b.Cursor.Row], b.Cursor.Col-1, b.Cursor.Col)
	b.Cursor.Col--
}

func (b *Buffer) DeleteLine() {
	if len(b.Content) == 1 {
		b.Content[0] = []byte{}
		return
	}
	b.Content = slices.Delete(b.Content, b.Cursor.Row, b.Cursor.Row+1)
	b.Cursor.Row--
}

func (b *Buffer) MoveCursorUp() {
	if b.Cursor.Row > 0 {
		b.Cursor.Row--
	}
	if b.Cursor.Col > len(b.Content[b.Cursor.Row])-1 {
		b.Cursor.Col = len(b.Content[b.Cursor.Row]) - 1
	}
	if len(b.Content[b.Cursor.Row]) == 0 {
		b.Cursor.Col = 0
	}
}
func (b *Buffer) MoveCursorDown() {

	if b.Cursor.Row < len(b.Content)-1 {
		b.Cursor.Row++
	}
	if b.Cursor.Col > len(b.Content[b.Cursor.Row])-1 {
		b.Cursor.Col = len(b.Content[b.Cursor.Row]) - 1
	}
	if len(b.Content[b.Cursor.Row]) == 0 {
		b.Cursor.Col = 0
	}

}
func (b *Buffer) MoveCursorLeft() {
	if b.Cursor.Col > 0 {
		b.Cursor.Col--
	}
}
func (b *Buffer) MoveCursorRight() {
	if b.Cursor.Col < len(b.Content[b.Cursor.Row]) {
		b.Cursor.Col++
	}
}
func (b *Buffer) SaveBuffer() {
}
func (b *Buffer) PrintBuffer() {
	terminal.ClearScreen()

	for i, line := range b.Content {
		for j, c := range line {
			if i == b.Cursor.Row && j == b.Cursor.Col {
				terminal.CursorColor()
				if b.Cursor.Type == InsertCursor {
					terminal.CursorBlinking()
				}
			}
			fmt.Printf("%c", c)
			terminal.ResetScreenAttributes()
		}
		if b.Cursor.Col == len(line) && i == b.Cursor.Row && len(line) > 0 {
			terminal.CursorColor()
			if b.Cursor.Type == InsertCursor {
				terminal.CursorBlinking()
			}
			fmt.Printf("%c", 32)
			terminal.ResetScreenAttributes()
		}
		if len(line) == 0 {
			if i == b.Cursor.Row {
				terminal.CursorColor()
				if b.Cursor.Type == InsertCursor {
					terminal.CursorBlinking()
				}
			}
			fmt.Printf("%c", 32)
			terminal.ResetScreenAttributes()
		}
		fmt.Printf("\r\n")
	}
	fmt.Println("Cursor Row: ", b.Cursor.Row, "Cursor Col: ", b.Cursor.Col)
}
