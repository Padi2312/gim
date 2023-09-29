package queue

type Queue interface {
	Enqueue(item interface{})
	Dequeue() interface{}
	IsEmpty() bool
}

type QueueImpl struct {
	items []interface{}
}

func NewQueueImpl() *QueueImpl {
	return &QueueImpl{}
}

func (q *QueueImpl) Enqueue(item interface{}) {
	q.items = append(q.items, item)
}

func (q *QueueImpl) Dequeue() interface{} {
	if q.IsEmpty() {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *QueueImpl) IsEmpty() bool {
	return len(q.items) == 0
}
