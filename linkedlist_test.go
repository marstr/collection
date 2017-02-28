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

func ExampleNewLinkedList() {
	subject1 := NewLinkedList('a', 'b', 'c', 'd', 'e')
	fmt.Println(subject1.Length())
	// Output: 5
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
	}

	if right == nil {
		t.Logf("unexpected nil value for right")
		t.Fail()
	} else if right.payload != 3 {
		t.Logf("got: %d\n want: %d", right.payload, 3)
		t.Fail()
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

func ExampleLinkedList_Sorta() {
	//subject := NewLinkedList("charlie", "alfa", "bravo", "delta")
	subject := NewLinkedList("foo", "bar")
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
	fmt.Println(subject.ToSlice())
	// Output: [2, 2, 3, 3, 6, 7]
}
