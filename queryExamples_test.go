package collection

import (
	"fmt"
	"sync"
)

func ExampleEnumerator_Count() {
	subject := AsEnumerable("str1", "str1", "str2")
	count1 := subject.Enumerate(nil).Count(func(a interface{}) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleEnumerator_CountAll() {
	subject := AsEnumerable('a', 'b', 'c', 'd', 'e')
	fmt.Println(subject.Enumerate(nil).CountAll())
	// Ouput: 5
}

func ExampleEnumerator_ElementAt() {
	done := make(chan struct{})
	defer close(done)
	fmt.Print(Fibonacci.Enumerate(done).ElementAt(4))
	// Output: 3
}

func ExampleFirst() {
	empty := NewQueue()
	notEmpty := NewQueue(1, 2, 3, 4)

	fmt.Println(First(empty))
	fmt.Println(First(notEmpty))
	// Output:
	// <nil> Enumerator encountered no elements
	// 1 <nil>
}

func ExampleLast() {
	subject := NewList(1, 2, 3, 4)
	fmt.Println(Last(subject))
	// Output: 4
}

func ExampleEnumerator_Last() {
	subject := AsEnumerable(1, 2, 3)
	fmt.Print(subject.Enumerate(nil).Last())
	//Output: 3
}

func ExampleMerge() {
	a := AsEnumerable(1, 2, 4)
	b := AsEnumerable(8, 16, 32)
	c := Merge(a, b)
	sum := 0
	for x := range c.Enumerate(nil) {
		sum += x.(int)
	}
	fmt.Println(sum)

	product := 1
	for y := range a.Enumerate(nil) {
		product *= y.(int)
	}
	fmt.Println(product)
	// Output:
	// 63
	// 8
}

func ExampleEnumerator_Reverse() {
	a := AsEnumerable(1, 2, 3).Enumerate(nil)
	a = a.Reverse()
	fmt.Println(a.ToSlice())
	// Output: [3 2 1]
}

func ExampleSelect() {
	const offset = 'a' - 1

	subject := AsEnumerable('a', 'b', 'c')
	subject = Select(subject, func(a interface{}) interface{} {
		return a.(rune) - offset
	})

	fmt.Println(ToSlice(subject))
	// Output: [1 2 3]
}

func ExampleEnumerator_Select() {
	subject := AsEnumerable('a', 'b', 'c').Enumerate(nil)
	const offset = 'a' - 1
	results := subject.Select(func(a interface{}) interface{} {
		return a.(rune) - offset
	})

	fmt.Println(results.ToSlice())
	// Output: [1 2 3]
}

func ExampleEnumerator_SelectMany() {

	type BrewHouse struct {
		Name  string
		Beers []interface{}
	}

	breweries := AsEnumerable(
		BrewHouse{
			"Mac & Jacks",
			[]interface{}{
				"African Amber",
				"Ibis IPA",
			},
		},
		BrewHouse{
			"Post Doc",
			[]interface{}{
				"Prereq Pale",
			},
		},
		BrewHouse{
			"Resonate",
			[]interface{}{
				"Comfortably Numb IPA",
				"Lithium Altbier",
			},
		},
		BrewHouse{
			"Triplehorn",
			[]interface{}{
				"Samson",
				"Pepper Belly",
			},
		},
	)

	beers := breweries.Enumerate(nil).SelectMany(func(brewer interface{}) Enumerator {
		return AsEnumerable(brewer.(BrewHouse).Beers...).Enumerate(nil)
	})

	for beer := range beers {
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
	done := make(chan struct{})
	defer close(done)

	trimmed := Take(Skip(Fibonacci, 1), 3)
	for entry := range trimmed.Enumerate(done) {
		fmt.Println(entry)
	}
	// Output:
	// 1
	// 1
	// 2
}

func ExampleEnumerator_Skip() {
	subject := AsEnumerable(1, 2, 3, 4, 5, 6, 7)
	skipped := subject.Enumerate(nil).Skip(5)
	for entry := range skipped {
		fmt.Println(entry)
	}
	// Output:
	// 6
	// 7
}

func ExampleTake() {
	done := make(chan struct{})
	defer close(done)

	taken := Take(Fibonacci, 4)
	for entry := range taken.Enumerate(done) {
		fmt.Println(entry)
	}
	// Output:
	// 0
	// 1
	// 1
	// 2
}

func ExampleEnumerator_Take() {
	done := make(chan struct{})
	defer close(done)

	taken := Fibonacci.Enumerate(done).Skip(4).Take(2)
	for entry := range taken {
		fmt.Println(entry)
	}
	// Output:
	// 3
	// 5
}

func ExampleTakeWhile() {
	taken := TakeWhile(Fibonacci, func(x interface{}, n uint) bool {
		return x.(int) < 10
	})
	for entry := range taken.Enumerate(nil) {
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
	taken := Fibonacci.Enumerate(nil).TakeWhile(func(x interface{}, n uint) bool {
		return x.(int) < 6
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
	base := AsEnumerable(1, 2, 4)
	left, right := base.Enumerate(nil).Tee()
	var wg sync.WaitGroup
	wg.Add(2)

	product := 1
	go func() {
		for x := range left {
			product *= x.(int)
		}
		wg.Done()
	}()

	sum := 0
	go func() {
		for x := range right {
			sum += x.(int)
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
	subject := NewStack(9, 'a', "str1")
	result := UCount(subject, func(a interface{}) bool {
		_, ok := a.(string)
		return ok
	})
	fmt.Println(result)
	// Output: 1
}

func ExampleEnumerator_UCount() {
	subject := AsEnumerable("str1", "str1", "str2")
	count1 := subject.Enumerate(nil).UCount(func(a interface{}) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleUCountAll() {
	subject := NewStack(8, 9, 10, 11)
	fmt.Println(UCountAll(subject))
	// Output: 4
}

func ExampleEnumerator_UCountAll() {
	subject := AsEnumerable('a', 2, "str1")
	fmt.Println(subject.Enumerate(nil).UCountAll())
	// Output: 3
}

func ExampleEnumerator_Where() {
	done := make(chan struct{})
	defer close(done)
	results := Fibonacci.Enumerate(done).Where(func(a interface{}) bool {
		return a.(int) > 8
	}).Take(3)
	fmt.Println(results.ToSlice())
	// Output: [13 21 34]
}

func ExampleWhere() {
	results := Where(AsEnumerable(1, 2, 3, 4, 5), func(a interface{}) bool {
		return a.(int) < 3
	})
	fmt.Println(ToSlice(results))
	// Output: [1 2]
}
