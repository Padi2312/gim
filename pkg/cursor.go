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

func (c *Cursor) Up() {
	c.Y--
}

func (c *Cursor) Down() {
	c.Y++
}

func (c *Cursor) Left() {
	c.X--
}

func (c *Cursor) Right() {
	c.X++
}
