package pkg

type Content struct {
	Buffer [][]rune
}

func NewContent() *Content {
	buffer := make([][]rune, 1)
	buffer[0] = make([]rune, 0)
	return &Content{Buffer: buffer}
}

func (c *Content) LineLength(lineNumber int) int {
	// Minus one because it comes from cursor
	return len(c.Buffer[lineNumber-1])
}

func (c *Content) TotalLines() int {
	return len(c.Buffer)
}

func (c *Content) InsertAtCursorV2(keyEvent KeyEvent, cursor *Cursor) {
	_, y := c.getBufferIndices(cursor)

	switch keyEvent.Char {
	case '\n':
		c.Buffer[y] = append(c.Buffer[y], '\n')
		c.Buffer = append(c.Buffer, make([]rune, 0))
	default:
		c.Buffer[y] = append(c.Buffer[y], keyEvent.Char)
	}
}

func (c *Content) getBufferIndices(cursor *Cursor) (int, int) {
	return cursor.X - 1, cursor.Y - 1
}
