package collection

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// LinkedList encapsulates a list where each entry is aware of only the next entry in the list.
type LinkedList struct {
	first  *llNode
	last   *llNode
	length uint
	key    sync.RWMutex
}

type llNode struct {
	payload interface{}
	next    *llNode
}

// Comparator is a function which evaluates two values to determine their relation to one another.
// - Zero is returned when `a` and `b` are equal.
// - Positive numbers are returned when `a` is greater than `b`.
// - Negative numbers are returned when `a` is less than `b`.
type Comparator func(a, b interface{}) (int, error)

// A collection of errors that may be thrown by functions in this file.
var (
	ErrListIsEmpty    = errors.New("targeted list is empty")
	ErrUnexpectedType = errors.New("value was of an unexpected type")
)

// NewLinkedList instantiates a new LinkedList with the entries provided.
func NewLinkedList(entries ...interface{}) *LinkedList {
	list := &LinkedList{}

	for _, entry := range entries {
		list.AddBack(entry)
	}

	return list
}

// AddBack creates an entry in the LinkedList that is logically at the back of the list.
func (list *LinkedList) AddBack(entry interface{}) {
	toAppend := &llNode{
		payload: entry,
	}

	list.key.Lock()
	defer list.key.Unlock()

	list.length++

	if list.first == nil {
		list.first = toAppend
		list.last = toAppend
		return
	}

	list.last.next = toAppend
	list.last = toAppend
}

// AddFront creates an entry in the LinkedList that is logically at the front of the list.
func (list *LinkedList) AddFront(entry interface{}) {
	toAppend := &llNode{
		payload: entry,
	}

	list.key.Lock()
	defer list.key.Unlock()

	list.length++

	if list.first == nil {
		list.first = toAppend
		list.last = toAppend
	}

	toAppend.next = list.first
	list.first = toAppend
}

// IsEmpty tests the list to determine if it is populate or not.
func (list *LinkedList) IsEmpty() bool {
	list.key.RLock()
	defer list.key.RUnlock()

	return list.first == nil
}

// Length returns the number of elements present in the LinkedList.
func (list *LinkedList) Length() uint {
	list.key.RLock()
	defer list.key.RUnlock()

	return list.length
}

// PeekBack returns the entry logicall stored at the back of the list without removing it.
func (list *LinkedList) PeekBack() (interface{}, error) {
	list.key.RLock()
	defer list.key.RUnlock()

	if list.last == nil {
		return nil, ErrListIsEmpty
	}
	return list.last.payload, nil
}

// PeekFront returns the entry logically stored at the front of this list without removing it.
func (list *LinkedList) PeekFront() (interface{}, error) {
	list.key.RLock()
	defer list.key.RUnlock()

	if list.first == nil {
		return nil, ErrListIsEmpty
	}
	return list.first.payload, nil
}

// RemoveFront returns the entry logically stored at the front of this list and removes it.
func (list *LinkedList) RemoveFront() (interface{}, error) {
	list.key.Lock()
	defer list.key.Unlock()

	if list.first == nil {
		return nil, ErrListIsEmpty
	}

	retval := list.first.payload

	list.first = list.first.next
	list.length--

	if 0 == list.length {
		list.last = nil
	}

	return retval, nil
}

// RemoveBack returns the entry logically stored at the back of this list and removes it.
func (list *LinkedList) RemoveBack() (interface{}, error) {
	list.key.Lock()
	defer list.key.Unlock()

	if list.last == nil {
		return nil, ErrListIsEmpty
	}

	retval := list.last.payload
	list.length--

	if list.length == 0 {
		list.first = nil
	} else {
		current := list.first
		endIndex := list.length - uint(1)
		for i := uint(0); i < endIndex; i++ {
			current = current.next
		}
		current.next = nil
	}
	return retval, nil
}

// Sort rearranges the positions of the entries in this list so that they are
// ascending.
func (list *LinkedList) Sort(comparator Comparator) error {
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

// Sorta rearranges the position of string entries in this list so that they
// are ascending.
func (list *LinkedList) Sorta() error {
	list.key.Lock()
	defer list.key.Unlock()

	var err error
	list.first, err = mergeSort(list.first, func(a, b interface{}) (int, error) {
		castA, ok := a.(string)
		if !ok {
			return 0, ErrUnexpectedType
		}
		castB, ok := b.(string)
		if !ok {
			return 0, ErrUnexpectedType
		}

		return strings.Compare(castA, castB), nil
	})
	list.last = findLast(list.first)
	return err
}

// Sorti rearranges the position of integer entries in this list so that they
// are ascending.
func (list *LinkedList) Sorti() error {
	list.key.Lock()
	defer list.key.Unlock()

	var err error
	list.first, err = mergeSort(list.first, func(a, b interface{}) (int, error) {
		castA, ok := a.(int)
		if !ok {
			return 0, ErrUnexpectedType
		}
		castB, ok := b.(int)
		if !ok {
			return 0, ErrUnexpectedType
		}

		return castA - castB, nil
	})
	if err != nil {
		return err
	}
	list.last = findLast(list.first)
	return err
}

// Swap switches the positions in which two values are stored in this list.
// x and y represent the indexes of the items that should be swapped.
func (list *LinkedList) Swap(x, y uint) error {
	list.key.Lock()
	defer list.key.Unlock()

	if x >= list.length {
		return fmt.Errorf("index out of bounds 'x', wanted less than %d got %d", list.length, x)
	}
	if y >= list.length {
		return fmt.Errorf("index out of bound 'y', wanted less than %d got %d", list.length, y)
	}

	current := list.first

	var xNode, yNode *llNode

	i := uint(0)
	for current != nil {
		if i == x {
			xNode = current
		}
		if i == y {
			yNode = current
		}
		if xNode != nil && yNode != nil {
			break
		}
		i++
	}

	temp := xNode.payload
	xNode.payload = yNode.payload
	yNode.payload = temp
	return nil
}

// ToSlice converts the contents of the LinkedList into a slice.
func (list *LinkedList) ToSlice() []interface{} {
	list.key.RLock()
	defer list.key.RUnlock()

	entryCount := int(list.length)
	retval := make([]interface{}, entryCount)

	current := list.first
	for i := 0; i < entryCount; i++ {
		retval[i] = current.payload
		current = current.next
	}
	return retval
}

func findLast(head *llNode) *llNode {
	if head == nil {
		return nil
	}
	current := head
	for current.next != nil {
		current = current.next
	}
	return current
}

// merge takes two sorted lists and merges them into one sorted list.
// Behavior is undefined when you pass a non-sorted list as `left` or `right`
func merge(left, right *llNode, comparator Comparator) (first *llNode, err error) {
	curLeft := left
	curRight := right

	var last *llNode

	appendResults := func(updated *llNode) {
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
			return
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

func mergeSort(head *llNode, comparator Comparator) (*llNode, error) {
	if head == nil {
		return nil, nil
	}

	left, right := split(head)

	repair := func(left, right *llNode) *llNode {
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
func split(head *llNode) (left, right *llNode) {
	left = head
	if head == nil || head.next == nil {
		return
	}
	right = head
	sprinter := head
	for sprinter != nil && sprinter.next != nil {
		right = right.next
		sprinter = sprinter.next.next
	}
	return
}
