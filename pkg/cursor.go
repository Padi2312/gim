package pkg

import "github.com/buger/goterm"

type Cursor struct {
	X int
	Y int
}

func NewCursor() *Cursor {
	goterm.MoveCursor(1, 1)
	goterm.Flush()
	return &Cursor{
		X: 1,
		Y: 1,
	}
}

func (c *Cursor) Update() {
	goterm.MoveCursor(c.X, c.Y)
	goterm.Flush()
}

func (c *Cursor) SetX(x int) {
	c.X = x
	c.Update()
}

func (c *Cursor) SetY(y int) {
	c.Y = y
	c.Update()
}

func (c *Cursor) JumpStartOfLine(lineNumber int) {
	c.X = 1
	c.Y = lineNumber
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
	c.Y++
	c.Update()
}

func (c *Cursor) Left() {
	if c.X > 1 {
		c.X--
		c.Update()
	}
}

func (c *Cursor) Right() {
	c.X++
	c.Update()
}
