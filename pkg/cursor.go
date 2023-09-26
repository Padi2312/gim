package pkg

import "github.com/buger/goterm"

type Cursor struct {
	X       int
	Y       int
	content *Content
}

func NewCursor(content *Content) *Cursor {
	goterm.MoveCursor(1, 1)
	goterm.Flush()
	return &Cursor{
		X:       1,
		Y:       1,
		content: content,
	}
}

func (c *Cursor) Update() {
	goterm.MoveCursor(c.X, c.Y)
	goterm.Flush()
}

func (c *Cursor) SetX(x int) {
	c.X = x
}

func (c *Cursor) SetY(y int) {
	c.Y = y
}

func (c *Cursor) JumpStartOfLine(lineNumber int) {
	c.X = 1
	c.Y = lineNumber
}
func (c *Cursor) GotoLineEnd() {
	c.X = c.content.LineLength(c.Y) + 1
	c.Update()
}
func (c *Cursor) JumpEndOfLine(lineNumber int) {
	c.X = 1
	c.Y = lineNumber
}
func (c *Cursor) Up() {
	if c.Y > 1 {
		c.Y--
		c.Update()
	}
}

func (c *Cursor) Down() {
	if c.Y+1 <= c.content.LineCount() {
		c.Y++
		c.Update()
	}
}

func (c *Cursor) Left() {
	if c.X > 1 {
		c.X--
		c.Update()
	}
}

func (c *Cursor) Right() {
	if c.X <= c.content.LineLength(c.Y) {
		c.X++
		c.Update()
	}
}
