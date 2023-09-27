package pkg

import "github.com/buger/goterm"

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

func (d *Display) Init() {
	goterm.Clear()
	goterm.MoveCursor(d.cursor.X, d.cursor.Y)
	goterm.Flush()
}

func (d *Display) DrawChar(char rune) {
	goterm.MoveCursor(d.cursor.X, d.cursor.Y)
	goterm.Print(string(char))
	goterm.Flush()
}

func (d *Display) Update(chaneQueue *ChangeQueue) {
	for !chaneQueue.IsEmpty() {
		change := chaneQueue.Dequeue()
		switch change.ChangeInst {
		case Remove:
			goterm.MoveCursor(change.Column, change.Line)
			goterm.Print(*change.Content)
		case Write:
			goterm.MoveCursor(change.Column, change.Line)
			goterm.Print(*change.Content)
		case Set:
			goterm.MoveCursor(change.Column, change.Line)
		}
	}
	goterm.Flush()
}
