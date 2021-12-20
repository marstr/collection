package collection

import (
	"context"
	"testing"
)

func TestLinkedList_findLast_empty(t *testing.T) {
	if result := findLast(nil); result != nil {
		t.Logf("got: %v\nwant: %v", result, nil)
	}
}

func TestLinkedList_merge(t *testing.T) {
	testCases := []struct {
		Left     *LinkedList
		Right    *LinkedList
		Expected []int
		Comp     Comparator
	}{
		{
			NewLinkedList(1, 3, 5),
			NewLinkedList(2, 4),
			[]int{1, 2, 3, 4, 5},
			UncheckedComparatori,
		},
		{
			NewLinkedList(1, 2, 3),
			NewLinkedList(),
			[]int{1, 2, 3},
			UncheckedComparatori,
		},
		{
			NewLinkedList(),
			NewLinkedList(1, 2, 3),
			[]int{1, 2, 3},
			UncheckedComparatori,
		},
		{
			NewLinkedList(),
			NewLinkedList(),
			[]int{},
			UncheckedComparatori,
		},
		{
			NewLinkedList(1),
			NewLinkedList(1),
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
			NewLinkedList(),
			[]int{3},
			UncheckedComparatori,
		},
		{
			NewLinkedList(),
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
					t.Logf("got: %d want: %d", cursor.payload.(int), tc.Expected[i])
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

func TestLinkedList_mergeSort_repair(t *testing.T) {
	testCases := []*LinkedList{
		NewLinkedList(1, 2, "str1", 4, 5, 6),
		NewLinkedList(1, 2, 3, "str1", 5, 6),
		NewLinkedList(1, 'a', 3, 4, 5, 6),
		NewLinkedList(1, 2, 3, 4, 5, uint(8)),
		NewLinkedList("alpha", 0),
		NewLinkedList(0, "kappa"),
	}

	for _, tc := range testCases {
		t.Run(tc.String(), func(t *testing.T) {
			originalLength := tc.Length()
			originalElements := tc.Enumerate(context.Background()).ToSlice()
			originalContents := tc.String()

			if err := tc.Sorti(); err != ErrUnexpectedType {
				t.Log("`Sorti() should have thrown ErrUnexpectedType")
				t.Fail()
			}

			t.Logf("Contents:\n\tOriginal:   \t%s\n\tPost Merge: \t%s", originalContents, tc.String())

			if newLength := tc.Length(); newLength != originalLength {
				t.Logf("Length changed. got: %d want: %d", newLength, originalLength)
				t.Fail()
			}

			remaining := tc.Enumerate(context.Background()).ToSlice()

			for _, desired := range originalElements {
				found := false
				for i, got := range remaining {
					if got == desired {
						remaining = append(remaining[:i], remaining[i+1:]...)
						found = true
						break
					}
				}

				if !found {
					t.Logf("couldn't find element: %v", desired)
					t.Fail()
				}
			}
		})
	}
}

func UncheckedComparatori(a, b interface{}) (int, error) {
	return a.(int) - b.(int), nil
}

func TestLinkedList_RemoveBack_single(t *testing.T) {
	subject := NewLinkedList(1)
	subject.RemoveBack()
	if subject.Length() != 0 {
		t.Fail()
	}
}

func TestLinkedList_Sorti(t *testing.T) {
	testCases := []struct {
		*LinkedList
		Expected []int
	}{
		{
			NewLinkedList(),
			[]int{},
		},
		{
			NewLinkedList(1, 2, 3, 4),
			[]int{1, 2, 3, 4},
		},
		{
			NewLinkedList(0, -1, 2, 8, 9),
			[]int{-1, 0, 2, 8, 9},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.String(), func(t *testing.T) {
			if err := tc.Sorti(); err != nil {
				t.Error(err)
			}

			sorted := tc.ToSlice()

			if countSorted, countExpected := len(sorted), len(tc.Expected); countSorted != countExpected {
				t.Logf("got: %d want: %d", countSorted, countExpected)
				t.FailNow()
			}

			for i, entry := range sorted {
				cast, ok := entry.(int)
				if !ok {
					t.Errorf("Element was not an int: %v", entry)
				}

				if cast != tc.Expected[i] {
					t.Logf("got: %d want: %d at: %d", cast, tc.Expected[i], i)
					t.Fail()
				}
			}
		})
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
	if !(result == nil && ok == false) {
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
			if first.(int) != 2 {
				t.Logf("got %d, want %d", first.(int), 2)
				t.Fail()
			}
		} else {
			t.Logf("no item at position 0!")
			t.Fail()
		}

		if second, ok := subject.Get(1); ok {
			if second.(int) != 3 {
				t.Logf("got %d, want %d", second.(int), 3)
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
			if first.(int) != 1 {
				t.Logf("got %d, want %d", first.(int), 1)
				t.Fail()
			}
		} else {
			t.Logf("no item at position 0!")
			t.Fail()
		}

		if second, ok := subject.Get(1); ok {
			if second.(int) != 2 {
				t.Logf("got %d, want %d", second.(int), 2)
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
			if first.(int) != 1 {
				t.Logf("got %d, want %d", first.(int), 1)
				t.Fail()
			}
		} else {
			t.Logf("no item at position 0!")
			t.Fail()
		}

		if second, ok := subject.Get(1); ok {
			if second.(int) != 3 {
				t.Logf("got %d, want %d", second.(int), 3)
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
