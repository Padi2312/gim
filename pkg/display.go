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

func (d *Display) DrawChar(char rune) {
	goterm.MoveCursor(d.cursor.X, d.cursor.Y)
	goterm.Print(string(char))
	goterm.Flush()
}

func (d *Display) Update(chaneQueue *ChangeQueue) {
	for !chaneQueue.IsEmpty() {
		change := chaneQueue.Dequeue()
		if change.ChangeInst == Remove {
			goterm.MoveCursor(change.From, change.Line)
			goterm.Print(change.Content)
		} else {
			goterm.MoveCursor(change.From, change.Line)
			goterm.Print(change.Content)
		}
	}
	goterm.Flush()
}
