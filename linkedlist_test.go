package collection

import "fmt"
import "testing"

func ExampleLinkedList_AddFront() {
	subject := NewLinkedList(2, 3)
	subject.AddFront(1)
	result, _ := subject.PeekFront()
	fmt.Println(result)
	// Output: 1
}

func ExampleLinkedList_AddBack() {
	subject := NewLinkedList(2, 3, 5)
	subject.AddBack(8)
	result, _ := subject.PeekBack()
	fmt.Println(result)
	fmt.Println(subject.Length())
	// Output:
	// 8
	// 4
}

func ExampleLinkedList_Get() {
	subject := NewLinkedList(2, 3, 5, 8)
	val, _ := subject.Get(2)
	fmt.Println(val)
	// Output: 5
}

func ExampleLinkedList_Get_OutsideBounds() {
	subject := NewLinkedList(2, 3, 5, 8, 13, 21)
	fmt.Print(subject.Get(10))
	// Output: <nil> false
}

func ExampleNewLinkedList() {
	subject1 := NewLinkedList('a', 'b', 'c', 'd', 'e')
	fmt.Println(subject1.Length())

	slice := []interface{}{1, 2, 3, 4, 5, 6}
	subject2 := NewLinkedList(slice...)
	fmt.Println(subject2.Length())
	// Output:
	// 5
	// 6
}

func ExampleLinkedList_Sort() {
	// Sorti sorts into ascending order, this example demonstrates sorting
	// into descending order.
	subject := NewLinkedList(2, 4, 3, 5, 7, 7)
	subject.Sort(func(a, b interface{}) (int, error) {
		castA, ok := a.(int)
		if !ok {
			return 0, ErrUnexpectedType
		}
		castB, ok := b.(int)
		if !ok {
			return 0, ErrUnexpectedType
		}

		return castB - castA, nil
	})
	fmt.Println(subject)
	// Output: [7 7 5 4 3 2]
}

func ExampleLinkedList_Sorta() {
	subject := NewLinkedList("charlie", "alfa", "bravo", "delta")
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
	subject := NewLinkedList(7, 3, 2, 2, 3, 6)
	subject.Sorti()
	fmt.Println(subject)
	// Output: [2 2 3 3 6 7]
}

func ExampleLinkedList_String() {
	subject := NewLinkedList(1, 2, 3)
	fmt.Println(subject)
	// Output: [1 2 3]
}

func ExampleLinkedList_String_Long() {
	subject := NewLinkedList()
	for i := 0; i < 20; i++ {
		subject.AddBack(i)
	}
	fmt.Println(subject)
	// Output: [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 ...]
}

func ExampleLinkedList_Swap() {
	subject := NewLinkedList(2, 3, 5, 8, 13)
	subject.Swap(1, 3)
	fmt.Println(subject)
	// Output: [2 8 5 3 13]
}

func TestSplit_Even(t *testing.T) {
	subject := NewLinkedList(1, 2, 3, 4)

	left, right := split(subject.first)
	if left == nil {
		t.Logf("unexpected nil value for left")
		t.Fail()
	} else if left.payload != 1 {
		t.Logf("got: %d\nwant: %d", left.payload, 1)
		t.Fail()
	}

	if right == nil {
		t.Logf("unexpected nil for right")
		t.Fail()
	} else if right.payload != 3 {
		t.Logf("got: %d\nwant: %d", right.payload, 3)
	}
}

func TestSplit_Odd(t *testing.T) {
	subject := NewLinkedList(1, 2, 3, 4, 5)

	left, right := split(subject.first)

	if left == nil {
		t.Logf("unexpected nil value for left")
		t.Fail()
	} else if left.payload != 1 {
		t.Logf("got: %d\n want: %d", left.payload, 1)
		t.Fail()
	} else if last := findLast(left).payload; last != 2 {
		t.Logf("got:\n%d\nwant:\n%d", last, 2)
		t.Fail()
	}

	if right == nil {
		t.Logf("unexpected nil value for right")
		t.Fail()
	} else if right.payload != 3 {
		t.Logf("got:\n%d\nwant:\n%d", right.payload, 3)
		t.Fail()
	} else if last := findLast(right).payload; last != 5 {
		t.Logf("got:\n%d\nwant:%d", last, 5)
	}
}

func TestSplit_Empty(t *testing.T) {
	subject := NewLinkedList()

	left, right := split(subject.first)

	if left != nil {
		t.Logf("got: %v\nwant: %v", left, nil)
		t.Fail()
	}

	if right != nil {
		t.Logf("got: %v\nwant: %v", right, nil)
		t.Fail()
	}
}

func TestSplit_Single(t *testing.T) {
	subject := NewLinkedList(1)

	left, right := split(subject.first)

	if left == nil {
		t.Logf("unexpected nil value for left")
		t.Fail()
	} else if left.payload != 1 {
		t.Logf("got: %d\nwant: %d", left.payload, 1)
		t.Fail()
	}

	if right != nil {
		t.Logf("got: %v\nwant: %v", right, nil)
		t.Fail()
	}

	if last := findLast(left).payload; last != 1 {
		t.Logf("got:\n%d\nwant:\n%d", last, 1)
		t.Fail()
	}
}

func TestSplit_Double(t *testing.T) {
	subject := NewLinkedList(1, 2)
	left, right := split(subject.first)

	if left == nil {
		t.Logf("unexpected nil value for left")
		t.Fail()
	} else if left.payload != 1 {
		t.Logf("got: %d\nwant: %d", left.payload, 1)
	}

	if right == nil {
		t.Logf("unexpected nil value for right")
		t.Fail()
	} else if right.payload != 2 {
		t.Logf("got: %d\nwant: %d", right.payload, 2)
	}
}

func UncheckedComparatori(a, b interface{}) (int, error) {
	return a.(int) - b.(int), nil
}

func TestMerge_EmptyLeft(t *testing.T) {
	left := NewLinkedList()

	expected := []interface{}{1, 2, 3}
	right := NewLinkedList(expected...)

	merged, err := merge(left.first, right.first, UncheckedComparatori)
	if err != nil {
		t.Error(err)
	}

	if len(expected) != 3 {
		// If it was manually changed, go ahead and modify the line above.
		// If it was programatically changed, something is wrong and you should investigate.
		t.Logf("expected has been modified")
		t.FailNow()
	}

	current := merged
	for i := 0; i < len(expected); i++ {
		if expected[i] != current.payload {
			t.Logf("got: %d\nwant: %d", current.payload, expected[i])
			t.Fail()
		}
		current = current.next
	}
}

func TestMerge_EmptyRight(t *testing.T) {
	right := NewLinkedList()

	expected := []interface{}{1, 2, 3}
	left := NewLinkedList(expected...)

	merged, err := merge(left.first, right.first, UncheckedComparatori)
	if err != nil {
		t.Error(err)
	}

	if len(expected) != 3 {
		// If it was manually changed, go ahead and modify the line above.
		// If it was programatically changed, something is wrong and you should investigate.
		t.Logf("expected has been modified")
		t.FailNow()
	}

	current := merged
	for i := 0; i < len(expected); i++ {
		if expected[i] != current.payload {
			t.Logf("got: %d\nwant: %d", current.payload, expected[i])
			t.Fail()
		}
		current = current.next
	}
}

func TestMerge_BothPopulated(t *testing.T) {
	expected := []interface{}{1, 2, 3, 4, 5}

	odds := NewLinkedList(1, 3, 5)
	evens := NewLinkedList(2, 4)

	merged, err := merge(odds.first, evens.first, UncheckedComparatori)
	if err != nil {
		t.Error(err)
	}

	if len(expected) != 5 {
		// If it was manually changed, go ahead and modify the line above.
		// If it was programatically changed, something is wrong and you should investigate.
		t.Logf("expected has been modified")
		t.FailNow()
	}

	current := merged
	for i := 0; i < len(expected); i++ {
		if expected[i] != current.payload {
			t.Logf("got: %d\nwant: %d", current.payload, expected[i])
			t.Fail()
		}
		current = current.next
	}
}

func TestMerge_BothEmpty(t *testing.T) {
	result, err := merge(nil, nil, UncheckedComparatori)
	if result != nil {
		t.Logf("got:\n%v\nwant:\n%v\n", result, nil)
		t.Fail()
	}
	if err != nil {
		t.Logf("got:\n%v\nwant:\n%v\n", result, nil)
		t.Fail()
	}
}
