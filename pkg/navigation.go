package pkg

type Navigation struct {
	content *Content
	cursor  *Cursor
}

func NewNavigation(content *Content, cursor *Cursor) *Navigation {
	return &Navigation{
		content: content,
		cursor:  cursor,
	}
}

func (n *Navigation) MoveLeft() {

}

func (n *Navigation) MoveRight(isInsert bool) {
	currentLineLength := n.content.LineLength(n.cursor.Y)
	if isInsert {
		currentLineLength++
	}

	if n.cursor.X <= currentLineLength {
		n.cursor.Right()
	}
}
