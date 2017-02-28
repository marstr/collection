package collection

// Enumerable exposes a new syntax for querying familiar data structures.
type Enumerable struct {
	output chan interface{}
}

// Any tests an Enumerable to see if there are any elements present.
func (iter *Enumerable) Any() bool {
	for range iter.output {
		return true
	}
	return false
}

// AsEnumerable allows for easy conversion to an Enumerable from a slice.
func AsEnumerable(entries ...interface{}) *Enumerable {
	retval := &Enumerable{
		output: make(chan interface{}),
	}

	go func() {
		for _, entry := range entries {
			retval.output <- entry
		}
		close(retval.output)
	}()

	return retval
}

// Predicate defines an interface for funcs that make some logical test.
type Predicate func(interface{}) bool

// AsEnumerablec allows for easy conversion to an Enumerable from a chan.
func AsEnumerablec(channel chan interface{}) *Enumerable {
	return &Enumerable{
		output: channel,
	}
}

// Count iterates over a list and keeps a running tally of the number of elements
// satisfy a predicate.
func (iter *Enumerable) Count(p Predicate) int {
	tally := 0
	for entry := range iter.output {
		if p(entry) {
			tally++
		}
	}
	return tally
}

// CountAll iterates over a list and keeps a running tally of how many it's seen.
func (iter *Enumerable) CountAll() int {
	tally := 0
	for range iter.output {
		tally++
	}
	return tally
}

// Select iterates over a list and returns a transformed item.
func (iter *Enumerable) Select(transform func(interface{}) interface{}) *Enumerable {
	retval := &Enumerable{
		output: make(chan interface{}),
	}

	go func() {
		for item := range iter.output {
			retval.output <- transform(item)
		}
		close(retval.output)
	}()

	return retval
}

// Tee splits athe results of a channel so that mulptiple actions can be taken on it.
func (iter *Enumerable) Tee() (*Enumerable, *Enumerable) {
	left := &Enumerable{
		output: make(chan interface{}),
	}
	right := &Enumerable{
		output: make(chan interface{}),
	}

	go func() {
		for entry := range iter.output {
			left.output <- entry
			right.output <- entry
		}
		close(left.output)
		close(right.output)
	}()

	return left, right
}

// ToChannel allows for conversion to use with traditional "range" syntax
func (iter *Enumerable) ToChannel() <-chan interface{} {
	return iter.output
}

// ToSlice places all iterated over values in a Slice for easy consumption.
func (iter *Enumerable) ToSlice() []interface{} {
	retval := make([]interface{}, 0)
	for entry := range iter.output {
		retval = append(retval, entry)
	}
	return retval
}

// Where iterates over a list and returns only the elements that satisfy a
// predicate.
func (iter *Enumerable) Where(predicate Predicate) *Enumerable {
	retval := &Enumerable{
		output: make(chan interface{}),
	}

	go func() {
		for item := range iter.output {
			if predicate(item) {
				retval.output <- item
			}
		}
		close(retval.output)
	}()

	return retval
}

// UCount iterates over a list and keeps a running tally of the number of elements
// satisfy a predicate.
func (iter *Enumerable) UCount(p Predicate) uint {
	tally := uint(0)
	for entry := range iter.output {
		if p(entry) {
			tally++
		}
	}
	return tally
}

// UCountAll iterates over a list and keeps a running tally of how many it's seen.
func (iter *Enumerable) UCountAll() uint {
	tally := uint(0)
	for range iter.output {
		tally++
	}
	return tally
}
