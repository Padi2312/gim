package pkg

import (
	"strings"

	"github.com/buger/goterm"
	"github.com/eiannone/keyboard"
)

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

func (c *Content) LineCount() int {
	return len(c.Buffer)
}

func (c *Content) InsertAtCursor(keyEvent KeyEvent, cursor *Cursor) {
	x, y := c.getBufferIndices(cursor)
	if keyEvent.Key == keyboard.KeyEnter {
		if c.LineLength(cursor.Y)+1 == cursor.X {
			c.Buffer[y] = append(c.Buffer[y], '\n')
			c.Buffer = append(c.Buffer, make([]rune, 0))
			// Call start to jump before and then navigate down
			cursor.JumpStartOfLine(cursor.Y)
			cursor.Down()
		} else {
			var appendingPart, remaining []rune
			appendingPart = make([]rune, len(c.Buffer[y][:x])) // Initialize appendingPart with the desired length
			remaining = make([]rune, len(c.Buffer[y][x:]))     // Initialize remaining with the desired length
			copy(appendingPart, c.Buffer[y][:x])
			copy(remaining, c.Buffer[y][x:])

			c.Buffer[y] = append(appendingPart, '\n')
			c.Buffer = append(c.Buffer, remaining)
			goterm.Print(strings.Repeat(" ", len(remaining)))
			goterm.Flush()

			cursor.JumpStartOfLine(cursor.Y)
			cursor.Down()
			goterm.Print(string(c.Buffer[cursor.Y-1]))
			goterm.Flush()
			cursor.GotoLineEnd()
		}
	} else if keyEvent.Key == keyboard.KeyBackspace {
		c.Buffer[y] = c.Buffer[y][:len(c.Buffer[y])-1]
		cursor.Left()
		goterm.Print(" ")
		goterm.Flush()
		// TODO: Figure out why i have to manually update the cursor here
		cursor.Update()
	} else if keyEvent.Key == 0 || keyEvent.Key == keyboard.KeySpace {
		if keyEvent.Key == keyboard.KeySpace {
			keyEvent.Char = ' '
		}

		if c.LineLength(cursor.Y)+1 == cursor.X {
			c.Buffer[y] = append(c.Buffer[y], keyEvent.Char)
		} else {
			var appendingPart, remaining []rune
			appendingPart = make([]rune, len(c.Buffer[y][:x])) // Initialize appendingPart with the desired length
			remaining = make([]rune, len(c.Buffer[y][x:]))     // Initialize remaining with the desired length
			copy(appendingPart, c.Buffer[y][:x])
			copy(remaining, c.Buffer[y][x:])
			appendingPart = append(appendingPart, keyEvent.Char)
			appendingPart = append(appendingPart, remaining...)
			c.Buffer[y] = appendingPart
		}
		goterm.Print(string(c.Buffer[y][x:]))
		goterm.Flush()

		cursor.Right()
	} else {
		return
	}
}

func (c *Content) getBufferIndices(cursor *Cursor) (int, int) {
	return cursor.X - 1, cursor.Y - 1
}
