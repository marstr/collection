package collection

import "fmt"

func ExampleNewStack_FromEmpty() {
	subject := NewStack()
	subject.Push("alfa")
	subject.Push("bravo")
	subject.Push("charlie")
	for !subject.IsEmpty() {
		val, _ := subject.Pop()
		fmt.Println(val)
	}
	// Output:
	// charlie
	// bravo
	// alfa
}

func ExampleNewStack_FromSlice() {
	subject := NewStack(1, 2, 3)
	for !subject.IsEmpty() {
		val, _ := subject.Pop()
		fmt.Println(val)
	}
	// Output:
	// 3
	// 2
	// 1
}
