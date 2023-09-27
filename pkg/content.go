package pkg

import "strings"

type Content struct {
	Buffer      [][]rune
	changeQueue *ChangeQueue
}

func NewContent(changeQueue *ChangeQueue) *Content {
	buffer := make([][]rune, 1)
	buffer[0] = make([]rune, 0)
	return &Content{
		Buffer:      buffer,
		changeQueue: changeQueue,
	}
}

func (c *Content) LineLength(lineNumber int) int {
	// Minus one because it comes from cursor
	return len(c.Buffer[lineNumber-1])
}

func (c *Content) TotalLines() int {
	return len(c.Buffer)
}

func (c *Content) DeleteBeforeCursor(cursor *Cursor) {
	_, y := c.getBufferIndices(cursor)
	c.Buffer[y] = c.Buffer[y][:len(c.Buffer[y])-1]
	nSpace := string(" ")
	changeRequestRemove := ChangeRequest{
		ChangeInst: Remove,
		Line:       cursor.Y,
		Column:     cursor.X - 1,
		Content:    &nSpace,
	}
	c.changeQueue.Enqueue(changeRequestRemove)
}

func (c *Content) InsertBeforeCursor(keyEvent KeyEvent, cursor *Cursor) {
	x, y := c.getBufferIndices(cursor)

	switch keyEvent.Char {
	case '\n':
		// Plus one because cursor is one field ahead of content in insert moe
		if cursor.X == c.LineLength(cursor.Y)+1 {
			c.Buffer = append(c.Buffer, make([]rune, 1))
		} else {
			remaining := c.Buffer[y][:x]
			newLineContent := c.Buffer[y][x:]
			valueRemove := string(strings.Repeat(" ", len(newLineContent)))
			changeRequestRemove := ChangeRequest{
				ChangeInst: Remove,
				Line:       cursor.Y,
				Column:     cursor.X,
				Content:    &valueRemove,
			}
			valueWrite := string(newLineContent)
			changeRequestWrite := ChangeRequest{
				ChangeInst: Write,
				Line:       cursor.Y + 1,
				Column:     1,
				Content:    &valueWrite,
			}
			c.changeQueue.Enqueue(changeRequestRemove)
			c.changeQueue.Enqueue(changeRequestWrite)

			c.Buffer[y] = remaining
			c.Buffer = append(c.Buffer, newLineContent)
		}
	default:
		if cursor.X == c.LineLength(cursor.Y)+1 {
			valueWrite := string(keyEvent.Char)
			changeRequestWrite := ChangeRequest{
				ChangeInst: Write,
				Line:       cursor.Y,
				Column:     cursor.X,
				Content:    &valueWrite,
			}
			c.changeQueue.Enqueue(changeRequestWrite)
			c.Buffer[y] = append(c.Buffer[y], keyEvent.Char)
		} else {
			originalLen := len(c.Buffer[y])
			newLine := make([]rune, originalLen+1)

			// Copy the existing content into the new slice
			copy(newLine, c.Buffer[y][:x])
			newLine[x] = keyEvent.Char
			copy(newLine[x+1:], c.Buffer[y][x:])

			// Create and enqueue the change requests
			value := string(keyEvent.Char)
			c.changeQueue.Enqueue(ChangeRequest{
				ChangeInst: Write,
				Line:       cursor.Y,
				Column:     cursor.X,
				Content:    &value,
			})

			remainingTextStr := string(c.Buffer[y][x:])
			c.changeQueue.Enqueue(ChangeRequest{
				ChangeInst: Write,
				Line:       cursor.Y,
				Column:     cursor.X + 1,
				Content:    &remainingTextStr,
			})

			// Update the buffer with the new slice
			c.Buffer[y] = newLine
		}

	}
}

func (c *Content) getBufferIndices(cursor *Cursor) (int, int) {
	return cursor.X - 1, cursor.Y - 1
}
