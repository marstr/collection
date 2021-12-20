package collection

import (
	"context"
	"sync"
)

// Stack implements a basic FILO structure.
type Stack[T any] struct {
	underlyer *LinkedList[T]
	key       sync.RWMutex
}

// NewStack instantiates a new FILO structure.
func NewStack[T any](entries ...T) *Stack[T] {
	retval := &Stack[T]{}
	retval.underlyer = NewLinkedList[T]()

	for _, entry := range entries {
		retval.Push(entry)
	}
	return retval
}

// Enumerate peeks at each element in the stack without mutating it.
func (stack *Stack[T]) Enumerate(ctx context.Context) Enumerator[T] {
	stack.key.RLock()
	defer stack.key.RUnlock()

	return stack.underlyer.Enumerate(ctx)
}

// IsEmpty tests the Stack to determine if it is populate or not.
func (stack *Stack[T]) IsEmpty() bool {
	stack.key.RLock()
	defer stack.key.RUnlock()
	return stack.underlyer == nil || stack.underlyer.IsEmpty()
}

// Push adds an entry to the top of the Stack.
func (stack *Stack[T]) Push(entry T) {
	stack.key.Lock()
	defer stack.key.Unlock()

	if nil == stack.underlyer {
		stack.underlyer = NewLinkedList[T]()
	}
	stack.underlyer.AddFront(entry)
}

// Pop returns the entry at the top of the Stack then removes it.
func (stack *Stack[T]) Pop() (T, bool) {
	stack.key.Lock()
	defer stack.key.Unlock()

	if nil == stack.underlyer {
		return *new(T), false
	}
	return stack.underlyer.RemoveFront()
}

// Peek returns the entry at the top of the Stack without removing it.
func (stack *Stack[T]) Peek() (T, bool) {
	stack.key.RLock()
	defer stack.key.RUnlock()
	return stack.underlyer.PeekFront()
}

// Size returns the number of entries populating the Stack.
func (stack *Stack[T]) Size() uint {
	stack.key.RLock()
	defer stack.key.RUnlock()
	if stack.underlyer == nil {
		return 0
	}
	return stack.underlyer.Length()
}
