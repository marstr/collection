package collection

import (
	"fmt"
	"testing"
)

func ExampleQueue_Add() {
	subject := &Queue{}
	subject.Add(1)
	subject.Add(2)
	res, _ := subject.Peek()
	fmt.Println(res)
	// Output: 1
}

func ExampleNewQueue() {
	empty := NewQueue()
	fmt.Println(empty.Length())

	populated := NewQueue(1, 2, 3, 5, 8, 13)
	fmt.Println(populated.Length())
	// Output:
	// 0
	// 6
}

func ExampleQueue_IsEmpty() {
	empty := NewQueue()
	fmt.Println(empty.IsEmpty())

	populated := NewQueue(1, 2, 3, 5, 8, 13)
	fmt.Println(populated.IsEmpty())
	// Output:
	// true
	// false
}

func ExampleQueue_Next() {
	subject := NewQueue(1, 2, 3, 5, 8, 13)
	for !subject.IsEmpty() {
		val, _ := subject.Next()
		fmt.Println(val)
	}
	// Output:
	// 1
	// 2
	// 3
	// 5
	// 8
	// 13
}

func TestQueue_Peek_DoesntRemove(t *testing.T) {
	expected := []interface{}{1, 2, 3}
	subject := NewQueue(expected...)
	if result, err := subject.Peek(); err != nil {
		t.Error(err)
	} else if result != expected[0] {
		t.Logf("got: %d\nwant: %d", result, 1)
		t.Fail()
	} else if count := subject.Length(); count != uint(len(expected)) {
		t.Logf("got: %d\nwant: %d", count, len(expected))
	}
}

func TestQueue_Length(t *testing.T) {
	empty := NewQueue()
	if count := empty.Length(); count != 0 {
		t.Logf("got: %d\nwant: %d", count, 0)
		t.Fail()
	}

	// Not the type magic number you're thinking of!
	// https://en.wikipedia.org/wiki/1729_(number)
	single := NewQueue(1729)
	if count := single.Length(); count != 1 {
		t.Logf("got: %d\nwant: %d", count, 1)
		t.Fail()
	}

	expectedMany := []interface{}{'a', 'b', 'c', 'd', 'e', 'e', 'f', 'g'}
	many := NewQueue(expectedMany...)
	if count := many.Length(); count != uint(len(expectedMany)) {
		t.Logf("got: %d\nwant: %d", count, len(expectedMany))
	}
}
