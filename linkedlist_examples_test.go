package collection_test

import (
	"fmt"

	"github.com/marstr/collection/v2"
)

func ExampleLinkedList_AddFront() {
	subject := collection.NewLinkedList(2, 3)
	subject.AddFront(1)
	result, _ := subject.PeekFront()
	fmt.Println(result)
	// Output: 1
}

func ExampleLinkedList_AddBack() {
	subject := collection.NewLinkedList(2, 3, 5)
	subject.AddBack(8)
	result, _ := subject.PeekBack()
	fmt.Println(result)
	fmt.Println(subject.Length())
	// Output:
	// 8
	// 4
}

func ExampleLinkedList_Enumerate() {
	subject := collection.NewLinkedList(2, 3, 5, 8)
	results := subject.Enumerate(nil).Select(func(a interface{}) interface{} {
		return -1 * a.(int)
	})
	for entry := range results {
		fmt.Println(entry)
	}
	// Output:
	// -2
	// -3
	// -5
	// -8
}

func ExampleLinkedList_Get() {
	subject := collection.NewLinkedList(2, 3, 5, 8)
	val, _ := subject.Get(2)
	fmt.Println(val)
	// Output: 5
}

func ExampleNewLinkedList() {
	subject1 := collection.NewLinkedList('a', 'b', 'c', 'd', 'e')
	fmt.Println(subject1.Length())

	slice := []interface{}{1, 2, 3, 4, 5, 6}
	subject2 := collection.NewLinkedList(slice...)
	fmt.Println(subject2.Length())
	// Output:
	// 5
	// 6
}

func ExampleLinkedList_Sort() {
	// Sorti sorts into ascending order, this example demonstrates sorting
	// into descending order.
	subject := collection.NewLinkedList(2, 4, 3, 5, 7, 7)
	subject.Sort(func(a, b interface{}) (int, error) {
		castA, ok := a.(int)
		if !ok {
			return 0, collection.ErrUnexpectedType
		}
		castB, ok := b.(int)
		if !ok {
			return 0, collection.ErrUnexpectedType
		}

		return castB - castA, nil
	})
	fmt.Println(subject)
	// Output: [7 7 5 4 3 2]
}

func ExampleLinkedList_Sorta() {
	subject := collection.NewLinkedList("charlie", "alfa", "bravo", "delta")
	subject.Sorta()
	for _, entry := range subject.ToSlice() {
		fmt.Println(entry.(string))
	}
	// Output:
	// alfa
	// bravo
	// charlie
	// delta
}

func ExampleLinkedList_Sorti() {
	subject := collection.NewLinkedList(7, 3, 2, 2, 3, 6)
	subject.Sorti()
	fmt.Println(subject)
	// Output: [2 2 3 3 6 7]
}

func ExampleLinkedList_String() {
	subject1 := collection.NewLinkedList()
	for i := 0; i < 20; i++ {
		subject1.AddBack(i)
	}
	fmt.Println(subject1)

	subject2 := collection.NewLinkedList(1, 2, 3)
	fmt.Println(subject2)
	// Output:
	// [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 ...]
	// [1 2 3]
}

func ExampleLinkedList_Swap() {
	subject := collection.NewLinkedList(2, 3, 5, 8, 13)
	subject.Swap(1, 3)
	fmt.Println(subject)
	// Output: [2 8 5 3 13]
}
