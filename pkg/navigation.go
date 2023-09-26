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

func (n *Navigation) MoveUp() {
	n.cursor.Up()
	if n.cursor.X > n.content.LineLength(n.cursor.Y) {
		n.cursor.SetX(n.content.LineLength(n.cursor.Y))
	}
}

func (n *Navigation) MoveDown() {
	if n.cursor.Y+1 <= n.content.TotalLines() {
		n.cursor.Down()
	}
}

func (n *Navigation) MoveDownLineBegin() {
	if n.cursor.Y+1 <= n.content.TotalLines() {
		n.cursor.Down()
		n.cursor.SetX(1)
	}
}

func (n *Navigation) MoveLeft() {
	n.cursor.Left()
}

func (n *Navigation) MoveRight(isInsert bool) {
	currentLineLength := n.content.LineLength(n.cursor.Y)
	if isInsert {
		currentLineLength++
	}

	if n.cursor.X+1 <= currentLineLength {
		n.cursor.Right()
	}
}
