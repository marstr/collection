package collection

import "fmt"

import "sync"

func ExampleEnumerable_Any() {
	empty := AsEnumerator()
	if empty.Any() {
		fmt.Println("Empty had some")
	} else {
		fmt.Println("Empty had none")
	}

	populated := AsEnumerator("str1")
	if populated.Any() {
		fmt.Println("Populated had some")
	} else {
		fmt.Println("Populated had none")
	}
	// Output:
	// Empty had none
	// Populated had some
}

func ExampleEnumerable_Count() {
	subject := AsEnumerator("str1", "str1", "str2")
	count1 := subject.Count(func(a interface{}) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleEnumerable_CountAll() {
	subject := AsEnumerator('a', 'b', 'c', 'd', 'e')
	fmt.Println(subject.CountAll())
	// Ouput: 5
}

func ExampleEnumerable_Select() {
	subject := AsEnumerator('a', 'b', 'c')
	const offset = 'a' - 1
	results := subject.Select(func(a interface{}) interface{} {
		return a.(rune) - offset
	})

	fmt.Println(results.ToSlice())
	// Output: [1 2 3]
}

func ExampleEnumerable_Tee() {
	base := AsEnumerator(1, 2, 4)
	left, right := base.Tee()

	var wg sync.WaitGroup

	var sumKey sync.Mutex
	sum := 0
	findSum := func(e <-chan interface{}) {
		defer wg.Done()
		for entry := range e {
			sumKey.Lock()
			sum += entry.(int)
			sumKey.Unlock()
		}
	}
	wg.Add(2)
	go findSum(left)
	go findSum(right)
	wg.Wait()

	fmt.Println(sum)
	//Output: 14
}

func ExampleEnumerable_UCount() {
	subject := AsEnumerator("str1", "str1", "str2")
	count1 := subject.UCount(func(a interface{}) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleEnumerable_UCountAll() {
	subject := AsEnumerator('a', 2, "str1")
	fmt.Println(subject.UCountAll())
	// Output: 3
}

func ExampleEnumerable_Where() {
	subject := AsEnumerator(1, 2, 3, 5, 8, 13, 21, 34)
	results := subject.Where(func(a interface{}) bool { return a.(int) > 8 }).ToSlice()
	fmt.Println(results)
	// Output: [13 21 34]
}
