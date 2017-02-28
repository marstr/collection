package collection

import (
	"sync"
)

// Queue implements a basic FIFO structure.
type Queue struct {
	underlyer *LinkedList
	key       sync.RWMutex
}

// Add places an item at the back of the Queue.
func (q *Queue) Add(entry interface{}) {
	q.key.Lock()
	defer q.key.Unlock()

	q.underlyer.AddBack(entry)
}

// Length returns the number of items in the Queue.
func (q *Queue) Length() uint {
	q.key.RLock()
	defer q.key.RUnlock()

	return q.underlyer.length
}

// Next removes and returns the next item in the Queue.
func (q *Queue) Next() (interface{}, error) {
	q.key.Lock()
	defer q.key.Unlock()

	return q.underlyer.RemoveFront()
}

// Peek returns the next item in the Queue without removing it.
func (q *Queue) Peek() (interface{}, error) {
	q.key.RLock()
	defer q.key.RUnlock()

	return q.underlyer.PeekFront()
}
