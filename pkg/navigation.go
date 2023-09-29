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
	if n.cursor.X >= n.content.LineLength(n.cursor.Y)-1 {
		n.cursor.SetX(n.content.LineLength(n.cursor.Y) - 1)
	}
}

func (n *Navigation) MoveDown() {
	if n.cursor.Y+1 < n.content.TotalLines() {
		n.cursor.Down()
		if n.cursor.X >= n.content.LineLength(n.cursor.Y)-1 {
			n.cursor.SetX(n.content.LineLength(n.cursor.Y) - 1)
		}
	}
}

func (n *Navigation) MoveDownLineBegin() {
	if n.cursor.Y+1 <= n.content.TotalLines() {
		n.cursor.Down()
		n.cursor.SetX(0)
	}
}

func (n *Navigation) JumpBeginOfLine() {
	n.cursor.SetX(0)
}

func (n *Navigation) JumpEndOfLine() {
	n.cursor.SetX(n.content.LineLength(n.cursor.Y) - 1)
}

func (n *Navigation) MoveLeft() {
	n.cursor.Left()
}

func (n *Navigation) MoveRight(isInsert bool) {
	currentLineLength := n.content.LineLength(n.cursor.Y)
	if isInsert {
		currentLineLength++
	}

	if n.cursor.X < currentLineLength-1 {
		n.cursor.Right()
	}
}

func (n *Navigation) Pos() (int, int) {
	return n.cursor.X, n.cursor.Y
}
