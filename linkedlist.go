package collection

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
)

// LinkedList encapsulates a list where each entry is aware of only the next entry in the list.
type LinkedList[T any] struct {
	first  *llNode[T]
	last   *llNode[T]
	length uint
	key    sync.RWMutex
}

type llNode[T any] struct {
	payload T
	next    *llNode[T]
	prev    *llNode[T]
}

// Comparator is a function which evaluates two values to determine their relation to one another.
// - Zero is returned when `a` and `b` are equal.
// - Positive numbers are returned when `a` is greater than `b`.
// - Negative numbers are returned when `a` is less than `b`.
type Comparator[T any] func(a, b T) (int, error)

// A collection of errors that may be thrown by functions in this file.
var (
	ErrUnexpectedType = errors.New("value was of an unexpected type")
)

// NewLinkedList instantiates a new LinkedList with the entries provided.
func NewLinkedList[T any](entries ...T) *LinkedList[T] {
	list := &LinkedList[T]{}

	for _, entry := range entries {
		list.AddBack(entry)
	}

	return list
}

// AddBack creates an entry in the LinkedList that is logically at the back of the list.
func (list *LinkedList[T]) AddBack(entry T) {
	list.key.Lock()
	defer list.key.Unlock()

	toAppend := &llNode[T]{
		payload: entry,
	}

	list.addNodeBack(toAppend)
}

func (list *LinkedList[T]) addNodeBack(node *llNode[T]) {

	list.length++

	node.prev = list.last

	if list.first == nil {
		list.first = node
		list.last = node
		return
	}

	list.last.next = node
	list.last = node
}

// AddFront creates an entry in the LinkedList that is logically at the front of the list.
func (list *LinkedList[T]) AddFront(entry T) {
	toAppend := &llNode[T]{
		payload: entry,
	}

	list.key.Lock()
	defer list.key.Unlock()

	list.addNodeFront(toAppend)
}

func (list *LinkedList[T]) addNodeFront(node *llNode[T]) {
	list.length++

	node.next = list.first
	if list.first == nil {
		list.last = node
	} else {
		list.first.prev = node
	}

	list.first = node
}

// Enumerate creates a new instance of Enumerable which can be executed on.
func (list *LinkedList[T]) Enumerate(ctx context.Context) Enumerator[T] {
	retval := make(chan T)

	go func() {
		list.key.RLock()
		defer list.key.RUnlock()
		defer close(retval)

		current := list.first
		for current != nil {
			select {
			case retval <- current.payload:
				// Intentionally Left Blank
			case <-ctx.Done():
				return
			}
			current = current.next
		}
	}()

	return retval
}

// Get finds the value from the LinkedList.
// pos is expressed as a zero-based index begining from the 'front' of the list.
func (list *LinkedList[T]) Get(pos uint) (T, bool) {
	list.key.RLock()
	defer list.key.RUnlock()
	node, ok := get(list.first, pos)
	if ok {
		return node.payload, true
	}
	return *new(T), false
}

// IsEmpty tests the list to determine if it is populate or not.
func (list *LinkedList[T]) IsEmpty() bool {
	list.key.RLock()
	defer list.key.RUnlock()

	return list.first == nil
}

// Length returns the number of elements present in the LinkedList.
func (list *LinkedList[T]) Length() uint {
	list.key.RLock()
	defer list.key.RUnlock()

	return list.length
}

// PeekBack returns the entry logicall stored at the back of the list without removing it.
func (list *LinkedList[T]) PeekBack() (T, bool) {
	list.key.RLock()
	defer list.key.RUnlock()

	if list.last == nil {
		return *new(T), false
	}
	return list.last.payload, true
}

// PeekFront returns the entry logically stored at the front of this list without removing it.
func (list *LinkedList[T]) PeekFront() (T, bool) {
	list.key.RLock()
	defer list.key.RUnlock()

	if list.first == nil {
		return *new(T), false
	}
	return list.first.payload, true
}

// RemoveFront returns the entry logically stored at the front of this list and removes it.
func (list *LinkedList[T]) RemoveFront() (T, bool) {
	list.key.Lock()
	defer list.key.Unlock()

	if list.first == nil {
		return *new(T), false
	}

	retval := list.first.payload

	list.first = list.first.next
	list.length--

	if list.length == 0 {
		list.last = nil
	}

	return retval, true
}

// RemoveBack returns the entry logically stored at the back of this list and removes it.
func (list *LinkedList[T]) RemoveBack() (T, bool) {
	list.key.Lock()
	defer list.key.Unlock()

	if list.last == nil {
		return *new(T), false
	}

	retval := list.last.payload
	list.length--

	if list.length == 0 {
		list.first = nil
	} else {
		list.last = list.last.prev
		list.last.next = nil
	}
	return retval, true
}

// removeNode
func (list *LinkedList[T]) removeNode(target *llNode[T]) {
	if target == nil {
		return
	}

	if target.next != nil {
		target.next.prev = target.prev
	}

	if target.prev != nil {
		target.prev.next = target.next
	}

	list.length--

	if list.length == 0 {
		list.first = nil
		list.last = nil
		return
	}

	if list.first == target {
		list.first = target.next
	}

	if list.last == target {
		list.last = target.prev
	}
}

// Sort rearranges the positions of the entries in this list so that they are
// ascending.
func (list *LinkedList[T]) Sort(comparator Comparator[T]) error {
	list.key.Lock()
	defer list.key.Unlock()
	var err error
	list.first, err = mergeSort(list.first, comparator)
	if err != nil {
		return err
	}
	list.last = findLast(list.first)
	return err
}

// String prints upto the first fifteen elements of the list in string format.
func (list *LinkedList[T]) String() string {
	list.key.RLock()
	defer list.key.RUnlock()

	builder := bytes.NewBufferString("[")
	current := list.first
	for i := 0; i < 15 && current != nil; i++ {
		builder.WriteString(fmt.Sprintf("%v ", current.payload))
		current = current.next
	}
	if current == nil || current.next == nil {
		builder.Truncate(builder.Len() - 1)
	} else {
		builder.WriteString("...")
	}
	builder.WriteRune(']')
	return builder.String()
}

// Swap switches the positions in which two values are stored in this list.
// x and y represent the indexes of the items that should be swapped.
func (list *LinkedList[T]) Swap(x, y uint) error {
	list.key.Lock()
	defer list.key.Unlock()

	var xNode, yNode *llNode[T]
	if temp, ok := get(list.first, x); ok {
		xNode = temp
	} else {
		return fmt.Errorf("index out of bounds 'x', wanted less than %d got %d", list.length, x)
	}
	if temp, ok := get(list.first, y); ok {
		yNode = temp
	} else {
		return fmt.Errorf("index out of bounds 'y', wanted less than %d got %d", list.length, y)
	}

	temp := xNode.payload
	xNode.payload = yNode.payload
	yNode.payload = temp
	return nil
}

// ToSlice converts the contents of the LinkedList into a slice.
func (list *LinkedList[T]) ToSlice() []T {
	return list.Enumerate(context.Background()).ToSlice()
}

func findLast[T any](head *llNode[T]) *llNode[T] {
	if head == nil {
		return nil
	}
	current := head
	for current.next != nil {
		current = current.next
	}
	return current
}

func get[T any](head *llNode[T], pos uint) (*llNode[T], bool) {
	for i := uint(0); i < pos; i++ {
		if head == nil {
			return nil, false
		}
		head = head.next
	}
	return head, true
}

// merge takes two sorted lists and merges them into one sorted list.
// Behavior is undefined when you pass a non-sorted list as `left` or `right`
func merge[T any](left, right *llNode[T], comparator Comparator[T]) (first *llNode[T], err error) {
	curLeft := left
	curRight := right

	var last *llNode[T]

	appendResults := func(updated *llNode[T]) {
		if last == nil {
			last = updated
		} else {
			last.next = updated
			last = last.next
		}
		if first == nil {
			first = last
		}
	}

	for curLeft != nil && curRight != nil {
		var res int
		if res, err = comparator(curLeft.payload, curRight.payload); nil != err {
			break // Don't return, stitch the remaining elements back on.
		} else if res < 0 {
			appendResults(curLeft)
			curLeft = curLeft.next
		} else {
			appendResults(curRight)
			curRight = curRight.next
		}
	}

	if curLeft != nil {
		appendResults(curLeft)
	}
	if curRight != nil {
		appendResults(curRight)
	}
	return
}

func mergeSort[T any](head *llNode[T], comparator Comparator[T]) (*llNode[T], error) {
	if head == nil {
		return nil, nil
	}

	left, right := split(head)

	repair := func(left, right *llNode[T]) *llNode[T] {
		lastLeft := findLast(left)
		lastLeft.next = right
		return left
	}

	var err error
	if left != nil && left.next != nil {
		left, err = mergeSort(left, comparator)
		if err != nil {
			return repair(left, right), err
		}
	}
	if right != nil && right.next != nil {
		right, err = mergeSort(right, comparator)
		if err != nil {
			return repair(left, right), err
		}
	}

	return merge(left, right, comparator)
}

// split breaks a list in half.
func split[T any](head *llNode[T]) (left, right *llNode[T]) {
	left = head
	if head == nil || head.next == nil {
		return
	}
	right = head
	sprinter := head
	prev := head
	for sprinter != nil && sprinter.next != nil {
		prev = right
		right = right.next
		sprinter = sprinter.next.next
	}
	prev.next = nil
	return
}
