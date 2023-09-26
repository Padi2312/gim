package pkg

type ChangeLog struct {
	ActionType string // "Insert", "Delete", "NewLine"
	Char       rune   // Character inserted or deleted
	X, Y       int    // Coordinates of the change
	LineLength int    // Line length BEFORE changes
}

type Content struct {
	Buffer    [][]rune
	Changelog []ChangeLog
}

func NewContent() *Content {
	buffer := make([][]rune, 1)
	buffer[0] = make([]rune, 0)
	return &Content{
		Buffer:    buffer,
		Changelog: make([]ChangeLog, 0),
	}
}

func (c *Content) LineLength(lineNumber int) int {
	// Minus one because it comes from cursor
	return len(c.Buffer[lineNumber-1])
}

func (c *Content) TotalLines() int {
	return len(c.Buffer)
}

func (c *Content) InsertAtCursorV2(keyEvent KeyEvent, cursor *Cursor) {
	x, y := c.getBufferIndices(cursor)

	switch keyEvent.Char {
	case '\n':
		if cursor.X == c.LineLength(cursor.Y)+1 {
			c.Buffer = append(c.Buffer, make([]rune, 0))
		} else {
			change := ChangeLog{"NewLine", keyEvent.Char, cursor.X, cursor.Y, len(c.Buffer[y])}
			c.Changelog = append(c.Changelog, change)

			remaining := c.Buffer[y][:x]
			newLineContent := c.Buffer[y][x:]

			c.Buffer[y] = remaining
			c.Buffer = append(c.Buffer, newLineContent)
		}
	default:
		change := ChangeLog{"Insert", keyEvent.Char, cursor.X, cursor.Y, len(c.Buffer[y])}
		c.Changelog = append(c.Changelog, change)
		c.Buffer[y] = append(c.Buffer[y], keyEvent.Char)
	}
}

func (c *Content) getBufferIndices(cursor *Cursor) (int, int) {
	return cursor.X - 1, cursor.Y - 1
}
