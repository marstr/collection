package collection

import (
	"sync"
)

// Stack implements a basic FILO structure.
type Stack struct {
	underlyer LinkedList
	key       sync.RWMutex
}

// Push adds an entry to the top of the Stack.
func (stack *Stack) Push(entry interface{}) {
	stack.key.Lock()
	defer stack.key.Unlock()
	stack.underlyer.AddFront(entry)
}

// Pop returns the entry at the top of the Stack then removes it.
func (stack *Stack) Pop() (interface{}, error) {
	stack.key.Lock()
	defer stack.key.Unlock()
	return stack.underlyer.RemoveFront()
}

// Peek returns the entry at the top of the Stack without removing it.
func (stack *Stack) Peek() (interface{}, error) {
	stack.key.RLock()
	defer stack.key.RUnlock()
	return stack.underlyer.PeekFront()
}
