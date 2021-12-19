package collection

import "testing"

func TestLinkedList_findLast_empty(t *testing.T) {
	if result := findLast[int](nil); result != nil {
		t.Logf("got: %v\nwant: %v", result, nil)
	}
}

func TestLinkedList_merge(t *testing.T) {
	testCases := []struct {
		Left     *LinkedList[int]
		Right    *LinkedList[int]
		Expected []int
		Comp     Comparator[int]
	}{
		{
			NewLinkedList[int](1, 3, 5),
			NewLinkedList[int](2, 4),
			[]int{1, 2, 3, 4, 5},
			UncheckedComparatori,
		},
		{
			NewLinkedList[int](1, 2, 3),
			NewLinkedList[int](),
			[]int{1, 2, 3},
			UncheckedComparatori,
		},
		{
			NewLinkedList[int](),
			NewLinkedList[int](1, 2, 3),
			[]int{1, 2, 3},
			UncheckedComparatori,
		},
		{
			NewLinkedList[int](),
			NewLinkedList[int](),
			[]int{},
			UncheckedComparatori,
		},
		{
			NewLinkedList[int](1),
			NewLinkedList[int](1),
			[]int{1, 1},
			UncheckedComparatori,
		},
		{
			NewLinkedList(2),
			NewLinkedList(1),
			[]int{1, 2},
			UncheckedComparatori,
		},
		{
			NewLinkedList(3),
			NewLinkedList[int](),
			[]int{3},
			UncheckedComparatori,
		},
		{
			NewLinkedList[int](),
			NewLinkedList(10),
			[]int{10},
			UncheckedComparatori,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			result, err := merge(tc.Left.first, tc.Right.first, tc.Comp)
			if err != nil {
				t.Error(err)
			}

			i := 0
			for cursor := result; cursor != nil; cursor, i = cursor.next, i+1 {
				if cursor.payload != tc.Expected[i] {
					t.Logf("got: %d want: %d", cursor.payload, tc.Expected[i])
					t.Fail()
				}
			}

			if expectedLength := len(tc.Expected); i != expectedLength {
				t.Logf("Unexpected length:\n\tgot: %d\n\twant: %d", i, expectedLength)
				t.Fail()
			}
		})
	}
}

func UncheckedComparatori(a, b int) (int, error) {
	return a - b, nil
}

func TestLinkedList_RemoveBack_single(t *testing.T) {
	subject := NewLinkedList(1)
	subject.RemoveBack()
	if subject.Length() != 0 {
		t.Fail()
	}
}

func TestLinkedList_split_Even(t *testing.T) {
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

func TestLinkedList_split_Odd(t *testing.T) {
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

func TestLinkedList_split_Empty(t *testing.T) {
	subject := NewLinkedList[*int]()

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

func TestLinkedList_split_Single(t *testing.T) {
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

func TestLinkedList_split_Double(t *testing.T) {
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

func TestLinkedList_Swap_OutOfBounds(t *testing.T) {
	subject := NewLinkedList(2, 3)
	if err := subject.Swap(0, 8); err == nil {
		t.Log("swap should have failed on y")
		t.Fail()
	}

	if err := subject.Swap(11, 1); err == nil {
		t.Logf("swap shoud have failed on x")
		t.Fail()
	}

	if count := subject.Length(); count != 2 {
		t.Logf("got: %d\nwant: %d", count, 2)
		t.Fail()
	}

	wantStr := "[2 3]"
	gotStr := subject.String()
	if wantStr != gotStr {
		t.Logf("got: %s\nwant: %s", gotStr, wantStr)
		t.Fail()
	}
}

func TestLinkedList_Get_OutsideBounds(t *testing.T) {
	subject := NewLinkedList(2, 3, 5, 8, 13, 21)
	result, ok := subject.Get(10)
	if !(result == 0 && ok == false) {
		t.Logf("got: %v %v\nwant: %v %v", result, ok, nil, false)
		t.Fail()
	}
}

func TestLinkedList_removeNode(t *testing.T) {
	removeHead := func(t *testing.T) {
		subject := NewLinkedList(1, 2, 3)

		subject.removeNode(subject.first)

		if subject.length != 2 {
			t.Logf("got %d, want %d", subject.length, 2)
			t.Fail()
		}

		if first, ok := subject.Get(0); ok {
			if first != 2 {
				t.Logf("got %d, want %d", first, 2)
				t.Fail()
			}
		} else {
			t.Logf("no item at position 0!")
			t.Fail()
		}

		if second, ok := subject.Get(1); ok {
			if second != 3 {
				t.Logf("got %d, want %d", second, 3)
				t.Fail()
			}
		} else {
			t.Logf("no item at position 1!")
			t.Fail()
		}
	}

	removeTail := func(t *testing.T) {
		subject := NewLinkedList(1, 2, 3)

		subject.removeNode(subject.last)

		if subject.length != 2 {
			t.Logf("got %d, want %d", subject.length, 2)
			t.Fail()
		}

		if first, ok := subject.Get(0); ok {
			if first != 1 {
				t.Logf("got %d, want %d", first, 1)
				t.Fail()
			}
		} else {
			t.Logf("no item at position 0!")
			t.Fail()
		}

		if second, ok := subject.Get(1); ok {
			if second != 2 {
				t.Logf("got %d, want %d", second, 2)
				t.Fail()
			}
		} else {
			t.Logf("no item at position 1!")
			t.Fail()
		}
	}

	removeMiddle := func(t *testing.T) {
		subject := NewLinkedList(1, 2, 3)

		subject.removeNode(subject.first.next)

		if subject.length != 2 {
			t.Logf("got %d, want %d", subject.length, 2)
			t.Fail()
		}

		if first, ok := subject.Get(0); ok {
			if first != 1 {
				t.Logf("got %d, want %d", first, 1)
				t.Fail()
			}
		} else {
			t.Logf("no item at position 0!")
			t.Fail()
		}

		if second, ok := subject.Get(1); ok {
			if second != 3 {
				t.Logf("got %d, want %d", second, 3)
				t.Fail()
			}
		} else {
			t.Logf("no item at position 1!")
			t.Fail()
		}
	}

	t.Run("RemoveHead", removeHead)
	t.Run("RemoveTail", removeTail)
	t.Run("RemoveMiddle", removeMiddle)
}
