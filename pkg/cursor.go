package pkg

type Cursor struct {
	X int
	Y int
}

func NewCursor() *Cursor {
	return &Cursor{
		X: 1,
		Y: 1,
	}
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

func (c *Cursor) Up() {
	if c.Y > 1 {
		c.Y--
	}
}

func (c *Cursor) Down() {
	c.Y++
}

func (c *Cursor) Left() {
	if c.X > 1 {
		c.X--
	}
}

func (c *Cursor) Right() {
	c.X++
}
