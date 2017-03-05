package collection

import (
	"errors"
	"sync"
)

// Enumerable offers a means of easily converting into a channel. It is most
// useful for types where mutability is not in question.
type Enumerable interface {
	Enumerate() <-chan interface{}
}

// Enumerator exposes a new syntax for querying familiar data structures.
type Enumerator <-chan interface{}

// Predicate defines an interface for funcs that make some logical test.
type Predicate func(interface{}) bool

var (
	errNoElements       = errors.New("Single.Enumerator encountered no elements")
	errMultipleElements = errors.New("Single.Enumerator encountered multiple elements")
)

// Any tests an Enumerator to see if there are any elements present.
func (iter Enumerator) Any() bool {
	for range iter {
		return true
	}
	return false
}

// AsEnumerator allows for easy conversion to an Enumerator from a slice.
func AsEnumerator(entries ...interface{}) Enumerator {
	retval := make(chan interface{})

	go func() {
		for _, entry := range entries {
			retval <- entry
		}
		close(retval)
	}()

	return retval
}

// Count iterates over a list and keeps a running tally of the number of elements
// satisfy a predicate.
func (iter Enumerator) Count(p Predicate) int {
	tally := 0
	for entry := range iter {
		if p(entry) {
			tally++
		}
	}
	return tally
}

// CountAll iterates over a list and keeps a running tally of how many it's seen.
func (iter Enumerator) CountAll() int {
	tally := 0
	for range iter {
		tally++
	}
	return tally
}

// Merge takes the results as it receives them from several channels and directs
// them into a single channel.
func Merge(channels ...<-chan interface{}) Enumerator {
	retval := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(len(channels))
	for _, item := range channels {
		go func(input <-chan interface{}) {
			defer wg.Done()
			for value := range input {
				retval <- value
			}
		}(item)
	}

	go func() {
		wg.Wait()
		close(retval)
	}()

	return retval
}

// Merge combines the results from this Enumerator with that of several others.
func (iter Enumerator) Merge(items ...<-chan interface{}) Enumerator {
	return Merge(append(items, iter)...)
}

// Reverse returns items in the opposite order it encountered them in.
func (iter Enumerator) Reverse() Enumerator {
	cache := NewStack()
	for entry := range iter {
		cache.Push(entry)
	}

	retval := make(chan interface{})

	go func() {
		for !cache.IsEmpty() {
			val, _ := cache.Pop()
			retval <- val
		}
		close(retval)
	}()
	return retval
}

// Select iterates over a list and returns a transformed item.
func (iter Enumerator) Select(transform func(interface{}) interface{}) Enumerator {
	retval := make(chan interface{})

	go func() {
		for item := range iter {
			retval <- transform(item)
		}
		close(retval)
	}()

	return retval
}

// Single retreives the only element from a list, or returns nil and an error.
func (iter Enumerator) Single() (interface{}, error) {
	var retval interface{}
	retError := errNoElements

	firstPass := true
	for entry := range iter {
		if firstPass {
			retval = entry
			retError = nil
		} else {
			retval = nil
			retError = errMultipleElements
			break
		}
		firstPass = false
	}
	return retval, retError
}

// Tee splits the results of a channel so that mulptiple actions can be taken on it.
func (iter Enumerator) Tee() (Enumerator, Enumerator) {
	left := make(chan interface{})
	right := make(chan interface{})

	go func() {
		for entry := range iter {
			left <- entry
			right <- entry
		}
		close(left)
		close(right)
	}()

	return left, right
}

// ToSlice places all iterated over values in a Slice for easy consumption.
func (iter Enumerator) ToSlice() []interface{} {
	retval := make([]interface{}, 0)
	for entry := range iter {
		retval = append(retval, entry)
	}
	return retval
}

// Where iterates over a list and returns only the elements that satisfy a
// predicate.
func (iter Enumerator) Where(predicate Predicate) Enumerator {
	retval := make(chan interface{})
	go func() {
		for item := range iter {
			if predicate(item) {
				retval <- item
			}
		}
		close(retval)
	}()

	return retval
}

// UCount iterates over a list and keeps a running tally of the number of elements
// satisfy a predicate.
func (iter Enumerator) UCount(p Predicate) uint {
	tally := uint(0)
	for entry := range iter {
		if p(entry) {
			tally++
		}
	}
	return tally
}

// UCountAll iterates over a list and keeps a running tally of how many it's seen.
func (iter Enumerator) UCountAll() uint {
	tally := uint(0)
	for range iter {
		tally++
	}
	return tally
}
