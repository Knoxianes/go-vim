package buffer

import "testing"

func TestBufferInsertChar(t *testing.T) {
	b := NewBuffer("testing_data")
	b.InsertChar('g')
	if b.Content[0][0] != 'g' {
		t.Errorf("InsertChar failed")
	}
}

func TestBufferInsertNewline(t *testing.T) {
	b := NewBuffer("testing_data")
	b.InsertNewline()
	if b.Content[1][0] != '\n' {
		t.Errorf("InsertNewline failed")
	}
}

func TestBufferDeleteChar(t *testing.T) {
	b := NewBuffer("testing_data")
	b.DeleteChar()
	if b.Content[0][0] != 'o' {
		t.Errorf("DeleteChar failed")
	}
	b.MoveCursorRight()
	b.DeleteChar()
	if b.Content[0][0] == 't' {
		t.Errorf("DeleteChar failed")
	}
}
func TestBufferDeleteLine(t *testing.T) {
	b := NewBuffer("testing_data")
	b.DeleteLine()
	if b.Content[0][0] == 'o' {
		t.Errorf("DeleteLine failed")
	}
}

func TestBufferCursorDown(t *testing.T) {
	b := NewBuffer("testing_data")
	b.MoveCursorDown()
	if b.Cursor.Row != 1 || b.Cursor.Col != 0 {
		t.Errorf("CursorDown failed")
	}
	for i := 0; i < 100; i++ {
		b.MoveCursorDown()
	}
	if b.Cursor.Row != (len(b.Content)-1) || b.Cursor.Col != 0 {
		t.Errorf("CursorDown failed")

	}
}

func TestBufferCursorUp(t *testing.T) {
	b := NewBuffer("testing_data")
	b.MoveCursorDown()
	b.MoveCursorUp()
	if b.Cursor.Row != 0 || b.Cursor.Col != 0 {
		t.Errorf("CursorUp failed")
	}

	b.MoveCursorUp()
	if b.Cursor.Row != 0 || b.Cursor.Col != 0 {
		t.Errorf("CursorUp failed")

	}
}

func TestBufferCursorRight(t *testing.T) {
	b := NewBuffer("testing_data")
	b.MoveCursorRight()
	if b.Cursor.Col != 1 || b.Cursor.Row != 0 {
		t.Errorf("CursorRight failed")
	}

	for i := 0; i < 100; i++ {
		b.MoveCursorRight()
	}
	if b.Cursor.Col != len(b.Content[0]) || b.Cursor.Row != 0 {
		t.Errorf("CursorRight failed")
	}
}

func TestBufferCursorLeft(t *testing.T) {
	b := NewBuffer("testing_data")
	b.MoveCursorRight()
	b.MoveCursorLeft()
	if b.Cursor.Col != 0 || b.Cursor.Row != 0 {
		t.Errorf("CursorLeft failed")
	}

	b.MoveCursorLeft()
	if b.Cursor.Col != 0 || b.Cursor.Row != 0 {
		t.Errorf("CursorLeft failed")
	}
}
