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

func (c *Content) InsertBeforeCursor(keyEvent KeyEvent, cursor *Cursor) {
	x, y := c.getBufferIndices(cursor)

	switch keyEvent.Char {
	case '\n':
		// Plus one because cursor is one field ahead of content in insert moe
		if cursor.X == c.LineLength(cursor.Y)+1 {
			c.Buffer = append(c.Buffer, make([]rune, 0))
		} else {
			remaining := c.Buffer[y][:x]
			newLineContent := c.Buffer[y][x:]
			changeRequestRemove := ChangeRequest{
				ChangeInst: Remove,
				Line:       cursor.Y,
				From:       cursor.X,
				To:         cursor.X + len(newLineContent),
				Content:    strings.Repeat(" ", len(newLineContent)),
			}
			changeRequestWrite := ChangeRequest{
				ChangeInst: Write,
				Line:       cursor.Y + 1,
				From:       1,
				To:         len(newLineContent),
				Content:    string(newLineContent),
			}
			c.changeQueue.Enqueue(changeRequestRemove)
			c.changeQueue.Enqueue(changeRequestWrite)

			c.Buffer[y] = remaining
			c.Buffer = append(c.Buffer, newLineContent)
		}
	default:
		if cursor.X == c.LineLength(cursor.Y)+1 {
			changeRequestWrite := ChangeRequest{
				ChangeInst: Write,
				Line:       cursor.Y,
				From:       cursor.X,
				To:         cursor.X,
				Content:    string(keyEvent.Char),
			}
			c.changeQueue.Enqueue(changeRequestWrite)
			c.Buffer[y] = append(c.Buffer[y], keyEvent.Char)
		} else {
			firstText := make([]rune, len(c.Buffer[y][:x]))
			remainingText := make([]rune, len(c.Buffer[y][x:]))

			copy(firstText, c.Buffer[y][:x])
			copy(remainingText, c.Buffer[y][x:])

			// Append character to start text
			firstText = append(firstText, keyEvent.Char)
			// Add remaining text back
			firstText = append(firstText, remainingText...)

			changeRequestWrite := ChangeRequest{
				ChangeInst: Write,
				Line:       cursor.Y,
				From:       cursor.X,
				To:         cursor.X,
				Content:    string(keyEvent.Char),
			}
			changeRequestWriteTwo := ChangeRequest{
				ChangeInst: Write,
				Line:       cursor.Y,
				From:       cursor.X + 1,
				To:         cursor.X,
				Content:    string(remainingText),
			}
			c.changeQueue.Enqueue(changeRequestWrite)
			c.changeQueue.Enqueue(changeRequestWriteTwo)
			c.Buffer[y] = firstText
		}

	}
}

func (c *Content) getBufferIndices(cursor *Cursor) (int, int) {
	return cursor.X - 1, cursor.Y - 1
}
