package collection

// Enumerable exposes a new syntax for querying familiar data structures.
type Enumerable <-chan interface{}

// Any tests an Enumerable to see if there are any elements present.
func (iter Enumerable) Any() bool {
	for range iter {
		return true
	}
	return false
}

// AsEnumerable allows for easy conversion to an Enumerable from a slice.
func AsEnumerable(entries ...interface{}) Enumerable {
	retval := make(chan interface{})

	go func() {
		for _, entry := range entries {
			retval <- entry
		}
		close(retval)
	}()

	return retval
}

// Predicate defines an interface for funcs that make some logical test.
type Predicate func(interface{}) bool

// Count iterates over a list and keeps a running tally of the number of elements
// satisfy a predicate.
func (iter Enumerable) Count(p Predicate) int {
	tally := 0
	for entry := range iter {
		if p(entry) {
			tally++
		}
	}
	return tally
}

// CountAll iterates over a list and keeps a running tally of how many it's seen.
func (iter Enumerable) CountAll() int {
	tally := 0
	for range iter {
		tally++
	}
	return tally
}

// Select iterates over a list and returns a transformed item.
func (iter Enumerable) Select(transform func(interface{}) interface{}) Enumerable {
	retval := make(chan interface{})

	go func() {
		for item := range iter {
			retval <- transform(item)
		}
		close(retval)
	}()

	return retval
}

// Tee splits athe results of a channel so that mulptiple actions can be taken on it.
func (iter Enumerable) Tee() (Enumerable, Enumerable) {
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
func (iter Enumerable) ToSlice() []interface{} {
	retval := make([]interface{}, 0)
	for entry := range iter {
		retval = append(retval, entry)
	}
	return retval
}

// Where iterates over a list and returns only the elements that satisfy a
// predicate.
func (iter Enumerable) Where(predicate Predicate) Enumerable {
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
func (iter Enumerable) UCount(p Predicate) uint {
	tally := uint(0)
	for entry := range iter {
		if p(entry) {
			tally++
		}
	}
	return tally
}

// UCountAll iterates over a list and keeps a running tally of how many it's seen.
func (iter Enumerable) UCountAll() uint {
	tally := uint(0)
	for range iter {
		tally++
	}
	return tally
}
