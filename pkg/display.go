package pkg

type Display struct {
	cursor *Cursor
}

func NewDisplay(cursor *Cursor) *Display {
	return &Display{
		cursor: cursor,
	}
}
