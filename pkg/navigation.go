package pkg

type Navigation struct {
	content     *Content
	cursor      *Cursor
	changeQueue *ChangeQueue
}

func NewNavigation(content *Content, cursor *Cursor, changeQueue *ChangeQueue) *Navigation {
	return &Navigation{
		content:     content,
		cursor:      cursor,
		changeQueue: changeQueue,
	}
}

func (n *Navigation) MoveUp() {
	n.cursor.Up()
	if n.cursor.X >= n.content.LineLength(n.cursor.Y) {
		n.cursor.SetX(n.content.LineLength(n.cursor.Y))
	}
	n.RequestChange()
}

func (n *Navigation) MoveDown() {
	if n.cursor.Y+1 <= n.content.TotalLines() {
		n.cursor.Down()
		if n.cursor.X >= n.content.LineLength(n.cursor.Y) {
			n.cursor.SetX(n.content.LineLength(n.cursor.Y))
		}
	}
	n.RequestChange()
}

func (n *Navigation) MoveDownLineBegin() {
	if n.cursor.Y+1 <= n.content.TotalLines() {
		n.cursor.Down()
		n.cursor.SetX(1)
		n.RequestChange()
	}
}

func (n *Navigation) MoveLeft() {
	n.cursor.Left()
	n.RequestChange()
}

func (n *Navigation) MoveRight(isInsert bool) {
	currentLineLength := n.content.LineLength(n.cursor.Y)
	if isInsert {
		currentLineLength++
	}

	if n.cursor.X+1 <= currentLineLength {
		n.cursor.Right()
		n.RequestChange()
	}
}

func (n *Navigation) RequestChange() {
	n.changeQueue.Enqueue(ChangeRequest{
		ChangeInst: Set,
		Line:       n.cursor.Y,
		Column:     n.cursor.X,
	})
}
