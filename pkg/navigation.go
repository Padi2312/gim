package pkg

import (
	"strings"
)

type Navigation struct {
	content *Content
	cursor  *Cursor
}

func NewNavigation(content *Content) *Navigation {
	return &Navigation{
		content: content,
		cursor:  NewCursor(),
	}
}

func (n *Navigation) MoveUp() {
	n.cursor.Up()
	if n.cursor.X > n.content.LineLength(n.cursor.Y) {
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

func (n *Navigation) WordForward() {
	currentLine := n.content.Buffer[n.cursor.Y][n.cursor.X:]
	vimWordBoundaries := "\t\n\v\f.,;:!?/ "

	for i := range currentLine {
		if strings.ContainsRune(vimWordBoundaries, currentLine[i]) {
			if i+1 <= n.content.LineLength(n.cursor.Y)-1 {
				// Plus one because we wont jump on the whitespace instead to the next char
				n.cursor.SetX(n.cursor.X + i + 1)
				return
			}
		}
	}

	// If cursor is not in last line we jump down one line to beginning
	if n.cursor.X == n.content.LineLength(n.cursor.Y)-1 && n.cursor.Y < n.content.TotalLines()-1 {
		n.MoveDown()
		n.JumpBeginOfLine()
	} else {
		// Otherwise we jump to end
		n.JumpEndOfLine()
	}
}

func (n *Navigation) WordBackward() {
	currentLine := n.content.Buffer[n.cursor.Y][:n.cursor.X]
	vimWordBoundaries := "\t\n\v\f.,;:!?/ "

	for i := range currentLine {
		i = len(currentLine) - 1 - i
		// If cursor is right before a whitespace ignore the char before the cursor
		if strings.ContainsRune(vimWordBoundaries, currentLine[i]) && i != n.cursor.X-1 {
			n.cursor.SetX(i + 1)
			return
		}
	}

	if n.cursor.X != 0 {
		n.JumpBeginOfLine()
	} else {
		n.MoveUp()
		n.JumpEndOfLine()
		n.WordBackward()
	}

}

func (n *Navigation) Pos() (int, int) {
	return n.cursor.X, n.cursor.Y
}
