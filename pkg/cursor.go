package pkg

type Cursor struct {
	X int
	Y int
}

const cursorLeftBorder = 0
const cursorTopBorder = 0

func NewCursor() *Cursor {
	return &Cursor{
		X: cursorLeftBorder,
		Y: cursorTopBorder,
	}
}

func (c *Cursor) SetX(x int) {
	if x < cursorLeftBorder {
		panic("Error on setting cursor X-value: Out of bounds")
	}
	c.X = x
}

func (c *Cursor) SetY(y int) {
	if y < cursorTopBorder{
		panic("Error on setting cursor X-value: Out of bounds")
	}
	c.Y = y
}

func (c *Cursor) JumpStartOfLine(lineNumber int) {
	c.X = cursorLeftBorder
	c.Y = lineNumber
}

func (c *Cursor) Up() {
	if c.Y > cursorTopBorder {
		c.Y--
	}
}

func (c *Cursor) Down() {
	c.Y++
}

func (c *Cursor) Left() {
	if c.X > cursorLeftBorder {
		c.X--
	}
}

func (c *Cursor) Right() {
	c.X++
}
