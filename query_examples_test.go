package collection_test

import (
	"context"
	"fmt"
	"sync"

	"github.com/marstr/collection/v2"
)

func ExampleEnumerableSlice_Enumerate() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// When a single value is provided, and it is an array or slice, each value in the array or slice is treated as an enumerable value.
	originalInts := []int{1, 2, 3, 4, 5}
	wrappedInts := collection.EnumerableSlice[int](originalInts)

	for entry := range wrappedInts.Enumerate(ctx) {
		fmt.Print(entry)
	}
	fmt.Println()

	// It's easy to convert arrays to slices for these enumerations as well.
	originalStrings := [7]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"}
	wrappedStrings := collection.EnumerableSlice[string](originalStrings[:])
	for entry := range wrappedStrings.Enumerate(ctx) {
		fmt.Println(entry)
	}
	// Output:
	// 12345
	// red
	// orange
	// yellow
	// green
	// blue
	// indigo
	// violet
}

func ExampleEnumerator_Count() {
	subject := collection.AsEnumerable("str1", "str1", "str2")
	count1 := subject.Enumerate(context.Background()).Count(func(a string) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleEnumerator_CountAll() {
	subject := collection.AsEnumerable('a', 'b', 'c', 'd', 'e')
	fmt.Println(subject.Enumerate(context.Background()).CountAll())
	// Output: 5
}

func ExampleEnumerator_ElementAt() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// ElementAt leaves the Enumerator open, creating a memory leak unless remediated,
	// context.Context should be cancelled to indicate that no further reads are coming.
	fmt.Print(collection.Fibonacci.Enumerate(ctx).ElementAt(4))
	// Output: 3
}

func ExampleFirst() {
	empty := collection.NewQueue[int]()
	notEmpty := collection.NewQueue(1, 2, 3, 4)

	fmt.Println(collection.First[int](empty))
	fmt.Println(collection.First[int](notEmpty))

	// Output:
	// 0 enumerator encountered no elements
	// 1 <nil>
}

func ExampleLast() {
	subject := collection.NewList(1, 2, 3, 4)
	fmt.Println(collection.Last[int](subject))
	// Output: 4
}

func ExampleEnumerator_Last() {
	subject := collection.AsEnumerable(1, 2, 3)
	fmt.Print(subject.Enumerate(context.Background()).Last())
	//Output: 3
}

func ExampleMerge() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := collection.AsEnumerable(1, 2, 4)
	b := collection.AsEnumerable(8, 16, 32)
	c := collection.Merge(a, b)
	sum := 0
	for x := range c.Enumerate(ctx) {
		sum += x
	}
	fmt.Println(sum)

	product := 1
	for y := range a.Enumerate(ctx) {
		product *= y
	}
	fmt.Println(product)
	// Output:
	// 63
	// 8
}

func ExampleEnumerator_Reverse() {
	a := collection.AsEnumerable(1, 2, 3).Enumerate(context.Background())
	a = a.Reverse()
	fmt.Println(a.ToSlice())
	// Output: [3 2 1]
}

func ExampleSelect() {
	const offset = 'a' - 1

	subject := collection.AsEnumerable[rune]('a', 'b', 'c')
	subject = collection.Select(subject, func(a rune) rune {
		return a - offset
	})

	fmt.Println(collection.ToSlice(subject))
	// Output: [1 2 3]
}

func ExampleEnumerator_SelectMany() {

	type BrewHouse struct {
		Name  string
		Beers collection.Enumerable[string]
	}

	breweries := collection.AsEnumerable(
		BrewHouse{
			"Mac & Jacks",
			collection.AsEnumerable[string](
				"African Amber",
				"Ibis IPA",
			),
		},
		BrewHouse{
			"Post Doc",
			collection.AsEnumerable[string](
				"Prereq Pale",
			),
		},
		BrewHouse{
			"Resonate",
			collection.AsEnumerable[string](
				"Comfortably Numb IPA",
				"Lithium Altbier",
			),
		},
		BrewHouse{
			"Triplehorn",
			collection.AsEnumerable[string](
				"Samson",
				"Pepper Belly",
			),
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	beers := collection.SelectMany(breweries, func(brewer BrewHouse) collection.Enumerator[string] {
		return brewer.Beers.Enumerate(ctx)
	})

	for beer := range beers.Enumerate(ctx) {
		fmt.Println(beer)
	}

	// Output:
	// African Amber
	// Ibis IPA
	// Prereq Pale
	// Comfortably Numb IPA
	// Lithium Altbier
	// Samson
	// Pepper Belly
}

func ExampleSkip() {
	trimmed := collection.Take(collection.Skip(collection.Fibonacci, 1), 3)
	for entry := range trimmed.Enumerate(context.Background()) {
		fmt.Println(entry)
	}
	// Output:
	// 1
	// 1
	// 2
}

func ExampleEnumerator_Skip() {
	subject := collection.AsEnumerable(1, 2, 3, 4, 5, 6, 7)
	skipped := subject.Enumerate(context.Background()).Skip(5)
	for entry := range skipped {
		fmt.Println(entry)
	}
	// Output:
	// 6
	// 7
}

func ExampleTake() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taken := collection.Take(collection.Fibonacci, 4)
	for entry := range taken.Enumerate(ctx) {
		fmt.Println(entry)
	}
	// Output:
	// 0
	// 1
	// 1
	// 2
}

func ExampleEnumerator_Take() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taken := collection.Fibonacci.Enumerate(ctx).Skip(4).Take(2)
	for entry := range taken {
		fmt.Println(entry)
	}
	// Output:
	// 3
	// 5
}

func ExampleTakeWhile() {
	taken := collection.TakeWhile(collection.Fibonacci, func(x, n uint) bool {
		return x < 10
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for entry := range taken.Enumerate(ctx) {
		fmt.Println(entry)
	}
	// Output:
	// 0
	// 1
	// 1
	// 2
	// 3
	// 5
	// 8
}

func ExampleEnumerator_TakeWhile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	taken := collection.Fibonacci.Enumerate(ctx).TakeWhile(func(x, n uint) bool {
		return x < 6
	})
	for entry := range taken {
		fmt.Println(entry)
	}
	// Output:
	// 0
	// 1
	// 1
	// 2
	// 3
	// 5
}

func ExampleEnumerator_Tee() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	base := collection.AsEnumerable(1, 2, 4)
	left, right := base.Enumerate(ctx).Tee()
	var wg sync.WaitGroup
	wg.Add(2)

	product := 1
	go func() {
		for x := range left {
			product *= x
		}
		wg.Done()
	}()

	sum := 0
	go func() {
		for x := range right {
			sum += x
		}
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("Sum: %d\n", sum)
	fmt.Printf("Product: %d\n", product)
	// Output:
	// Sum: 7
	// Product: 8
}

func ExampleUCount() {
	subject := collection.NewStack[any](9, 'a', "str1")
	result := collection.UCount[any](subject, func(a interface{}) bool {
		_, ok := a.(string)
		return ok
	})
	fmt.Println(result)
	// Output: 1
}

func ExampleEnumerator_UCount() {
	subject := collection.EnumerableSlice[string]([]string{"str1", "str1", "str2"})
	count1 := subject.Enumerate(context.Background()).UCount(func(a string) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleUCountAll() {
	subject := collection.NewStack(8, 9, 10, 11)
	fmt.Println(collection.UCountAll[int](subject))
	// Output: 4
}

func ExampleEnumerator_UCountAll() {
	subject := collection.EnumerableSlice[any]([]interface{}{'a', 2, "str1"})
	fmt.Println(subject.Enumerate(context.Background()).UCountAll())
	// Output: 3
}

func ExampleEnumerator_Where() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := collection.Fibonacci.Enumerate(ctx).Where(func(a uint) bool {
		return a > 8
	}).Take(3)
	fmt.Println(results.ToSlice())
	// Output: [13 21 34]
}

func ExampleWhere() {
	nums := collection.EnumerableSlice[int]([]int{1, 2, 3, 4, 5})
	results := collection.Where[int](nums, func(a int) bool {
		return a < 3
	})
	fmt.Println(collection.ToSlice(results))
	// Output: [1 2]
}
