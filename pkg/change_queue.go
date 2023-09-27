package pkg

type ChangeInst = int

const (
	Write ChangeInst = iota
	Remove
	Set
)

type ChangeRequest struct {
	ChangeInst ChangeInst
	Line       int
	Column     int
	Content    *string
}

type ChangeQueue struct {
	changes []ChangeRequest
}

func NewChangeQueue() *ChangeQueue {
	return &ChangeQueue{
		changes: make([]ChangeRequest, 0),
	}
}

func (c *ChangeQueue) Enqueue(changeRequest ChangeRequest) {
	c.changes = append(c.changes, changeRequest)
}

func (c *ChangeQueue) Dequeue() *ChangeRequest {
	if c.IsEmpty() {
		return nil
	} else {
		currentChange := c.changes[0]
		c.changes = c.changes[1:]
		return &currentChange
	}
}

func (c *ChangeQueue) IsEmpty() bool {
	return len(c.changes) == 0
}
