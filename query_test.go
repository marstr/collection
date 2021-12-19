package collection

import (
	"testing"
	"time"
)

func Test_Empty(t *testing.T) {
	if Any(Empty[int]()) {
		t.Log("empty should not have any elements")
		t.Fail()
	}

	if CountAll(Empty[int]()) != 0 {
		t.Log("empty should have counted to zero elements")
		t.Fail()
	}

	alwaysTrue := func(x int) bool {
		return true
	}

	if Count(Empty[int](), alwaysTrue) != 0 {
		t.Log("empty should have counted to zero even when discriminating")
		t.Fail()
	}
}

func BenchmarkEnumerator_Sum(b *testing.B) {
	var nums EnumerableSlice[int] = getInitializedSequentialArray[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slowNums := Select[int, int](nums, sleepIdentity[int])
		for range slowNums.Enumerate(nil) {
			// Intentionally Left Blank
		}
	}
}

func sleepIdentity[T any](val T) T {
	time.Sleep(2 * time.Millisecond)
	return val
}

func getInitializedSequentialArray[T ~int]() []T {
	rawNums := make([]T, 1000, 1000)
	for i := range rawNums {
		rawNums[i] = T(i + 1)
	}
	return rawNums
}
