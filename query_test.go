package collection

import "fmt"

func ExampleEnumerable_Count() {
	subject := AsEnumerable("str1", "str1", "str2")
	count1 := subject.Count(func(a interface{}) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleEnumerable_CountAll() {
	subject := AsEnumerable('a', 'b', 'c', 'd', 'e')
	fmt.Println(subject.CountAll())
	// Ouput: 5
}

func ExampleEnumerable_Select() {
	subject := AsEnumerable('a', 'b', 'c')
	const offset = 'a' - 1
	results := subject.Select(func(a interface{}) interface{} { return a.(rune) - offset }).ToSlice()
	fmt.Println(results)
	// Output: [1 2 3]
}

func ExampleEnumerable_UCount() {
	subject := AsEnumerable("str1", "str1", "str2")
	count1 := subject.UCount(func(a interface{}) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleEnumerable_UCountAll() {
	subject := AsEnumerable('a', 2, "str1")
	fmt.Println(subject.UCountAll())
	// Output: 3
}

func ExampleEnumerable_Where() {
	subject := AsEnumerable(1, 2, 3, 5, 8, 13, 21, 34)
	results := subject.Where(func(a interface{}) bool { return a.(int) > 8 }).ToSlice()
	fmt.Println(results)
	// Output: [13 21 34]
}
