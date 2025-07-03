package collection

import (
	"context"
	"errors"
	"runtime"
	"sync"
)

// Enumerable offers a means of easily converting into a channel. It is most
// useful for types where mutability is not in question.
type Enumerable[T any] interface {
	Enumerate(ctx context.Context) Enumerator[T]
}

// Enumerator exposes a new syntax for querying familiar data structures.
type Enumerator[T any] <-chan T

// Predicate defines an interface for funcs that make some logical test.
type Predicate[T any] func(T) bool

// Transform defines a function which takes a value, and returns some value based on the original.
type Transform[T any, E any] func(T) E

// Unfolder defines a function which takes a single value, and exposes many of them as an Enumerator
type Unfolder[T any, E any] func(T) Enumerator[E]

type emptyEnumerable[T any] struct{}

var (
	errNoElements       = errors.New("enumerator encountered no elements")
	errMultipleElements = errors.New("enumerator encountered multiple elements")
)

// IsErrorNoElements determines whethr or not the given error is the result of no values being
// returned when one or more were expected.
func IsErrorNoElements(err error) bool {
	return err == errNoElements
}

// IsErrorMultipleElements determines whether or not the given error is the result of multiple values
// being returned when one or zero were expected.
func IsErrorMultipleElements(err error) bool {
	return err == errMultipleElements
}

// Identity returns a trivial Transform which applies no operation on the value.
func Identity[T any]() Transform[T, T] {
	return func(value T) T {
		return value
	}
}

func Empty[T any]() Enumerable[T] {
	return &emptyEnumerable[T]{}
}

func (e emptyEnumerable[T]) Enumerate(ctx context.Context) Enumerator[T] {
	results := make(chan T)
	close(results)
	return results
}

// All tests whether or not all items present in an Enumerable meet a criteria.
func All[T any](subject Enumerable[T], p Predicate[T]) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return subject.Enumerate(ctx).All(p)
}

// All tests whether or not all items present meet a criteria.
func (iter Enumerator[T]) All(p Predicate[T]) bool {
	for entry := range iter {
		if !p(entry) {
			return false
		}
	}
	return true
}

// Any tests an Enumerable to see if there are any elements present.
func Any[T any](iterator Enumerable[T]) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for range iterator.Enumerate(ctx) {
		return true
	}
	return false
}

// Anyp tests an Enumerable to see if there are any elements present that meet a criteria.
func Anyp[T any](iterator Enumerable[T], p Predicate[T]) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for element := range iterator.Enumerate(ctx) {
		if p(element) {
			return true
		}
	}
	return false
}

type EnumerableSlice[T any] []T

func (f EnumerableSlice[T]) Enumerate(ctx context.Context) Enumerator[T] {
	results := make(chan T)

	go func() {
		defer close(results)
		for _, entry := range f {
			select {
			case results <- entry:
				// Intentionally Left Blank
			case <-ctx.Done():
				return
			}
		}
	}()

	return results
}

// AsEnumerable allows for easy conversion of a slice to a re-usable Enumerable object.
func AsEnumerable[T any](entries ...T) Enumerable[T] {
	return EnumerableSlice[T](entries)
}

// AsEnumerable stores the results of an Enumerator so the results can be enumerated over repeatedly.
func (iter Enumerator[T]) AsEnumerable() Enumerable[T] {
	return EnumerableSlice[T](iter.ToSlice())
}

// Count iterates over a list and keeps a running tally of the number of elements which satisfy a predicate.
func Count[T any](iter Enumerable[T], p Predicate[T]) int {
	return iter.Enumerate(context.Background()).Count(p)
}

// Count iterates over a list and keeps a running tally of the number of elements
// satisfy a predicate.
func (iter Enumerator[T]) Count(p Predicate[T]) int {
	tally := 0
	for entry := range iter {
		if p(entry) {
			tally++
		}
	}
	return tally
}

// CountAll iterates over a list and keeps a running tally of how many it's seen.
func CountAll[T any](iter Enumerable[T]) int {
	return iter.Enumerate(context.Background()).CountAll()
}

// CountAll iterates over a list and keeps a running tally of how many it's seen.
func (iter Enumerator[T]) CountAll() int {
	tally := 0
	for range iter {
		tally++
	}
	return tally
}

// Discard reads an enumerator to the end but does nothing with it.
// This method should be used in circumstances when it doesn't make sense to explicitly cancel the Enumeration.
func (iter Enumerator[T]) Discard() {
	for range iter {
		// Intentionally Left Blank
	}
}

// ElementAt retreives an item at a particular position in an Enumerator.
func ElementAt[T any](iter Enumerable[T], n uint) T {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return iter.Enumerate(ctx).ElementAt(n)
}

// ElementAt retreives an item at a particular position in an Enumerator.
func (iter Enumerator[T]) ElementAt(n uint) T {
	for i := uint(0); i < n; i++ {
		<-iter
	}
	return <-iter
}

// First retrieves just the first item in the list, or returns an error if there are no elements in the array.
func First[T any](subject Enumerable[T]) (retval T, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = errNoElements
	var isOpen bool

	if retval, isOpen = <-subject.Enumerate(ctx); isOpen {
		err = nil
	}

	return
}

// Last retreives the item logically behind all other elements in the list.
func Last[T any](iter Enumerable[T]) T {
	return iter.Enumerate(context.Background()).Last()
}

// Last retreives the item logically behind all other elements in the list.
func (iter Enumerator[T]) Last() (retval T) {
	for retval = range iter {
		// Intentionally Left Blank
	}
	return
}

type merger[T any] struct {
	originals []Enumerable[T]
}

func (m merger[T]) Enumerate(ctx context.Context) Enumerator[T] {
	retval := make(chan T)

	var wg sync.WaitGroup
	wg.Add(len(m.originals))
	for _, item := range m.originals {
		go func(input Enumerable[T]) {
			defer wg.Done()
			for value := range input.Enumerate(ctx) {
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

// Merge takes the results as it receives them from several channels and directs
// them into a single channel.
func Merge[T any](channels ...Enumerable[T]) Enumerable[T] {
	return merger[T]{
		originals: channels,
	}
}

// Merge takes the results of this Enumerator and others, and funnels them into
// a single Enumerator. The order of in which they will be combined is non-deterministic.
func (iter Enumerator[T]) Merge(others ...Enumerator[T]) Enumerator[T] {
	retval := make(chan T)

	var wg sync.WaitGroup
	wg.Add(len(others) + 1)

	funnel := func(prevResult Enumerator[T]) {
		for entry := range prevResult {
			retval <- entry
		}
		wg.Done()
	}

	go funnel(iter)
	for _, item := range others {
		go funnel(item)
	}

	go func() {
		wg.Wait()
		close(retval)
	}()
	return retval
}

type parallelSelecter[T any, E any] struct {
	original  Enumerable[T]
	operation Transform[T, E]
}

func (ps parallelSelecter[T, E]) Enumerate(ctx context.Context) Enumerator[E] {
	iter := ps.original.Enumerate(ctx)
	if cpus := runtime.NumCPU(); cpus != 1 {
		intermediate := splitN(iter, ps.operation, uint(cpus))
		return intermediate[0].Merge(intermediate[1:]...)
	}

	return Select(ps.original, ps.operation).Enumerate(ctx)
}

// ParallelSelect creates an Enumerable which will use all logically available CPUs to
// execute a Transform.
func ParallelSelect[T any, E any](original Enumerable[T], operation Transform[T, E]) Enumerable[E] {
	return parallelSelecter[T, E]{
		original:  original,
		operation: operation,
	}
}

// ParallelSelect will execute a Transform across all logical CPUs available to the current process.
//
// This is commented out, because Go 1.18 adds support for generics, but disallows methods from having type parameters
// not declared by their receivers.
//
//func (iter Enumerator[T]) ParallelSelect[E any](operation Transform[T, E]) Enumerator[E] {
//	if cpus := runtime.NumCPU(); cpus != 1 {
//		intermediate := iter.splitN(operation, uint(cpus))
//		return intermediate[0].Merge(intermediate[1:]...)
//	}
//	return iter
//}

type reverser[T any] struct {
	original Enumerable[T]
}

// Reverse will enumerate all values of an enumerable, store them in a Stack, then replay them all.
func Reverse[T any](original Enumerable[T]) Enumerable[T] {
	return reverser[T]{
		original: original,
	}
}

func (r reverser[T]) Enumerate(ctx context.Context) Enumerator[T] {
	return r.original.Enumerate(ctx).Reverse()
}

// Reverse returns items in the opposite order it encountered them in.
func (iter Enumerator[T]) Reverse() Enumerator[T] {
	cache := NewStack[T]()
	for entry := range iter {
		cache.Push(entry)
	}

	retval := make(chan T)

	go func() {
		for !cache.IsEmpty() {
			val, _ := cache.Pop()
			retval <- val
		}
		close(retval)
	}()
	return retval
}

type selecter[T any, E any] struct {
	original  Enumerable[T]
	transform Transform[T, E]
}

func (s selecter[T, E]) Enumerate(ctx context.Context) Enumerator[E] {
	retval := make(chan E)

	go func() {
		defer close(retval)

		for item := range s.original.Enumerate(ctx) {
			select {
			case retval <- s.transform(item):
				// Intentionally Left Blank
			case <-ctx.Done():
				return
			}
		}
	}()

	return retval
}

// Select creates a reusable stream of transformed values.
func Select[T any, E any](subject Enumerable[T], transform Transform[T, E]) Enumerable[E] {
	return selecter[T, E]{
		original:  subject,
		transform: transform,
	}
}

// Select iterates over a list and returns a transformed item.
//
// This is commented out because Go 1.18 added support for
//
//func (iter Enumerator[T]) Select[E any](transform Transform[T, E]) Enumerator[E] {
//	retval := make(chan interface{})
//
//	go func() {
//		for item := range iter {
//			retval <- transform(item)
//		}
//		close(retval)
//	}()
//
//	return retval
//}

type selectManyer[T any, E any] struct {
	original Enumerable[T]
	toMany   Unfolder[T, E]
}

func (s selectManyer[T, E]) Enumerate(ctx context.Context) Enumerator[E] {
	retval := make(chan E)

	go func() {
		for parent := range s.original.Enumerate(ctx) {
			for child := range s.toMany(parent) {
				retval <- child
			}
		}
		close(retval)
	}()
	return retval
}

// SelectMany allows for unfolding of values.
func SelectMany[T any, E any](subject Enumerable[T], toMany Unfolder[T, E]) Enumerable[E] {
	return selectManyer[T, E]{
		original: subject,
		toMany:   toMany,
	}
}

//// SelectMany allows for flattening of data structures.
//func (iter Enumerator[T]) SelectMany[E any](lister Unfolder[T, E]) Enumerator[E] {
//	retval := make(chan E)
//
//	go func() {
//		for parent := range iter {
//			for child := range lister(parent) {
//				retval <- child
//			}
//		}
//		close(retval)
//	}()
//
//	return retval
//}

// Single retreives the only element from a list, or returns nil and an error.
func Single[T any](iter Enumerable[T]) (retval T, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = errNoElements

	firstPass := true
	for entry := range iter.Enumerate(ctx) {
		if firstPass {
			retval = entry
			err = nil
		} else {
			retval = *new(T)
			err = errMultipleElements
			break
		}
		firstPass = false
	}
	return
}

// Singlep retrieces the only element from a list that matches a criteria. If
// no match is found, or two or more are found, `Singlep` returns nil and an
// error.
func Singlep[T any](iter Enumerable[T], pred Predicate[T]) (retval T, err error) {
	iter = Where(iter, pred)
	return Single(iter)
}

type skipper[T any] struct {
	original  Enumerable[T]
	skipCount uint
}

func (s skipper[T]) Enumerate(ctx context.Context) Enumerator[T] {
	return s.original.Enumerate(ctx).Skip(s.skipCount)
}

// Skip creates a reusable stream which will skip the first `n` elements before iterating
// over the rest of the elements in an Enumerable.
func Skip[T any](subject Enumerable[T], n uint) Enumerable[T] {
	return skipper[T]{
		original:  subject,
		skipCount: n,
	}
}

// Skip retreives all elements after the first 'n' elements.
func (iter Enumerator[T]) Skip(n uint) Enumerator[T] {
	results := make(chan T)

	go func() {
		defer close(results)

		i := uint(0)
		for entry := range iter {
			if i < n {
				i++
				continue
			}
			results <- entry
		}
	}()

	return results
}

// splitN creates N Enumerators, each will be a subset of the original Enumerator and will have
// distinct populations from one another.
func splitN[T any, E any](iter Enumerator[T], operation Transform[T, E], n uint) []Enumerator[E] {
	results, cast := make([]chan E, n), make([]Enumerator[E], n)

	for i := uint(0); i < n; i++ {
		results[i] = make(chan E)
		cast[i] = results[i]
	}

	go func() {
		for i := uint(0); i < n; i++ {
			go func(addr uint) {
				defer close(results[addr])
				for {
					read, ok := <-iter
					if !ok {
						return
					}
					results[addr] <- operation(read)
				}
			}(i)
		}
	}()

	return cast
}

type taker[T any] struct {
	original Enumerable[T]
	n        uint
}

func (t taker[T]) Enumerate(ctx context.Context) Enumerator[T] {
	return t.original.Enumerate(ctx).Take(t.n)
}

// Take retreives just the first `n` elements from an Enumerable.
func Take[T any](subject Enumerable[T], n uint) Enumerable[T] {
	return taker[T]{
		original: subject,
		n:        n,
	}
}

// Take retreives just the first 'n' elements from an Enumerator.
func (iter Enumerator[T]) Take(n uint) Enumerator[T] {
	results := make(chan T)

	go func() {
		defer close(results)
		i := uint(0)
		for entry := range iter {
			if i >= n {
				return
			}
			i++
			results <- entry
		}
	}()

	return results
}

type takeWhiler[T any] struct {
	original Enumerable[T]
	criteria func(T, uint) bool
}

func (tw takeWhiler[T]) Enumerate(ctx context.Context) Enumerator[T] {
	return tw.original.Enumerate(ctx).TakeWhile(tw.criteria)
}

// TakeWhile creates a reusable stream which will halt once some criteria is no longer met.
func TakeWhile[T any](subject Enumerable[T], criteria func(T, uint) bool) Enumerable[T] {
	return takeWhiler[T]{
		original: subject,
		criteria: criteria,
	}
}

// TakeWhile continues returning items as long as 'criteria' holds true.
func (iter Enumerator[T]) TakeWhile(criteria func(T, uint) bool) Enumerator[T] {
	results := make(chan T)

	go func() {
		defer close(results)
		i := uint(0)
		for entry := range iter {
			if !criteria(entry, i) {
				return
			}
			i++
			results <- entry
		}
	}()

	return results
}

// Tee creates two Enumerators which will have identical contents as one another.
func (iter Enumerator[T]) Tee() (Enumerator[T], Enumerator[T]) {
	left, right := make(chan T), make(chan T)

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
func ToSlice[T any](iter Enumerable[T]) []T {
	return iter.Enumerate(context.Background()).ToSlice()
}

// ToSlice places all iterated over values in a Slice for easy consumption.
func (iter Enumerator[T]) ToSlice() []T {
	retval := make([]T, 0)
	for entry := range iter {
		retval = append(retval, entry)
	}
	return retval
}

type wherer[T any] struct {
	original Enumerable[T]
	filter   Predicate[T]
}

func (w wherer[T]) Enumerate(ctx context.Context) Enumerator[T] {
	retval := make(chan T)

	go func() {
		defer close(retval)
		for entry := range w.original.Enumerate(ctx) {
			if w.filter(entry) {
				retval <- entry
			}
		}
	}()

	return retval
}

// Where creates a reusable means of filtering a stream.
func Where[T any](original Enumerable[T], p Predicate[T]) Enumerable[T] {
	return wherer[T]{
		original: original,
		filter:   p,
	}
}

// Where iterates over a list and returns only the elements that satisfy a
// predicate.
func (iter Enumerator[T]) Where(predicate Predicate[T]) Enumerator[T] {
	retval := make(chan T)
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
func UCount[T any](iter Enumerable[T], p Predicate[T]) uint {
	return iter.Enumerate(context.Background()).UCount(p)
}

// UCount iterates over a list and keeps a running tally of the number of elements
// satisfy a predicate.
func (iter Enumerator[T]) UCount(p Predicate[T]) uint {
	tally := uint(0)
	for entry := range iter {
		if p(entry) {
			tally++
		}
	}
	return tally
}

// UCountAll iterates over a list and keeps a running tally of how many it's seen.
func UCountAll[T any](iter Enumerable[T]) uint {
	return iter.Enumerate(context.Background()).UCountAll()
}

// UCountAll iterates over a list and keeps a running tally of how many it's seen.
func (iter Enumerator[T]) UCountAll() uint {
	tally := uint(0)
	for range iter {
		tally++
	}
	return tally
}
