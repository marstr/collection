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
	fmt.Print(Fibonacci.Enumerate(nil).ElementAt(4))
	// Output: 3
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
	// Output: 63
}

func ExampleEnumerator_Reverse() {
	a := AsEnumerable(1, 2, 3)
	fmt.Println(a.Enumerate(nil).Reverse().ToSlice())
	// Output: [3 2 1]
}

func ExampleEnumerator_Select() {
	subject := AsEnumerable('a', 'b', 'c')
	const offset = 'a' - 1
	results := subject.Enumerate(nil).Select(func(a interface{}) interface{} {
		return a.(rune) - offset
	})

	fmt.Println(results.ToSlice())
	// Output: [1 2 3]
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

func ExampleEnumerator_Split() {
	a := AsEnumerable(1, 2, 4, 8, 16)
	left, right := a.Enumerate(nil).Split(Identity)

	var wg sync.WaitGroup
	wg.Add(2)

	leftSum := 0
	go func() {
		for x := range left {
			leftSum += x.(int)
		}
		wg.Done()
	}()

	rightSum := 0
	go func() {
		for y := range right {
			rightSum += y.(int)
		}
		wg.Done()
	}()
	wg.Wait()

	fmt.Print(leftSum + rightSum)
	// Output: 31
}

func ExampleEnumerator_Take() {
	done := make(chan struct{})
	defer close(done)

	subject := AsEnumerable(1, 2, 3, 4, 5, 6)
	taken := subject.Enumerate(done).Skip(2).Take(2)
	for entry := range taken {
		fmt.Println(entry)
	}
	// Output:
	// 3
	// 4
}

func ExampleEnumerator_TakeWhile() {
	done := make(chan struct{})
	defer close(done)

	taken := Fibonacci.Enumerate(done).TakeWhile(func(x interface{}, n uint) bool {
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

func ExampleEnumerator_UCount() {
	subject := AsEnumerable("str1", "str1", "str2")
	count1 := subject.Enumerate(nil).UCount(func(a interface{}) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleEnumerator_UCountAll() {
	subject := AsEnumerable('a', 2, "str1")
	fmt.Println(subject.Enumerate(nil).UCountAll())
	// Output: 3
}

func ExampleEnumerator_Where() {
	fibonnaci := AsEnumerable(1, 2, 3, 5, 8, 13, 21, 34)
	results := fibonnaci.Enumerate(nil).Where(func(a interface{}) bool {
		return a.(int) > 8
	})
	fmt.Println(results.ToSlice())
	// Output: [13 21 34]
}

func ExampleLast() {
	result := Last(AsEnumerable('a', 'b', 'c', 'd'))
	fmt.Println(result == 'd')
	// Output: true
}

func ExampleWhere() {
	results := Where(AsEnumerable(1, 2, 3, 4, 5), func(a interface{}) bool {
		return a.(int) < 3
	})
	fmt.Println(ToSlice(results))
	// Output: [1 2]
}
