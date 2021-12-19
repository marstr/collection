package collection

import (
	"sync"
)

// Queue implements a basic FIFO structure.
type Queue[T any] struct {
	underlyer *LinkedList[T]
	key       sync.RWMutex
}

// NewQueue instantiates a new FIFO structure.
func NewQueue[T any](entries ...T) *Queue[T] {
	retval := &Queue[T]{
		underlyer: NewLinkedList[T](entries...),
	}
	return retval
}

// Add places an item at the back of the Queue.
func (q *Queue[T]) Add(entry T) {
	q.key.Lock()
	defer q.key.Unlock()
	if nil == q.underlyer {
		q.underlyer = NewLinkedList[T]()
	}
	q.underlyer.AddBack(entry)
}

// Enumerate peeks at each element of this queue without mutating it.
func (q *Queue[T]) Enumerate(cancel <-chan struct{}) Enumerator[T] {
	q.key.RLock()
	defer q.key.RUnlock()
	return q.underlyer.Enumerate(cancel)
}

// IsEmpty tests the Queue to determine if it is populate or not.
func (q *Queue[T]) IsEmpty() bool {
	q.key.RLock()
	defer q.key.RUnlock()
	return q.underlyer == nil || q.underlyer.IsEmpty()
}

// Length returns the number of items in the Queue.
func (q *Queue[T]) Length() uint {
	q.key.RLock()
	defer q.key.RUnlock()
	if nil == q.underlyer {
		return 0
	}
	return q.underlyer.length
}

// Next removes and returns the next item in the Queue.
func (q *Queue[T]) Next() (T, bool) {
	q.key.Lock()
	defer q.key.Unlock()
	if q.underlyer == nil {
		return *new(T), false
	}
	return q.underlyer.RemoveFront()
}

// Peek returns the next item in the Queue without removing it.
func (q *Queue[T]) Peek() (T, bool) {
	q.key.RLock()
	defer q.key.RUnlock()
	if q.underlyer == nil {
		return *new(T), false
	}
	return q.underlyer.PeekFront()
}

// ToSlice converts a Queue into a slice.
func (q *Queue[T]) ToSlice() []T {
	q.key.RLock()
	defer q.key.RUnlock()

	if q.underlyer == nil {
		return []T{}
	}
	return q.underlyer.ToSlice()
}
