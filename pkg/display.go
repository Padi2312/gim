package pkg

import (
	"strings"

	"github.com/buger/goterm"
)

type Display struct {
	cursor  *Cursor
	content *Content
}

func NewDisplay(cursor *Cursor, content *Content) *Display {
	return &Display{
		cursor:  cursor,
		content: content,
	}
}

func (d *Display) DrawChar(char rune) {
	goterm.MoveCursor(d.cursor.X, d.cursor.Y)
	goterm.Print(string(char))
	goterm.Flush()
}

func (d *Display) DrawChanges(changes []ChangeLog) {
	for _, change := range changes {
		switch change.ActionType {
		case "Insert":
			goterm.MoveCursor(change.X, change.Y)
			goterm.Print(string(change.Char))

		case "NewLine":
			goterm.MoveCursor(change.X, change.Y)

			// Clear from the cursor to the end of the line by printing spaces
			// The number of spaces would be the difference between the end of the line and the cursor
			nSpaces := change.LineLength - d.content.LineLength(change.Y)
			strings.Repeat(" ", nSpaces)
			for i := 0; i < nSpaces; i++ {
				goterm.Print(" ")
			}
		}
	}

	// Clear the changes slice
	changes = changes[:0]

	// Update the terminal screen
	goterm.Flush()
}
