package collection

import (
	"testing"
	"time"
)

func BenchmarkEnumerator_Sum(b *testing.B) {
	nums := AsEnumerable(getInitializedSequentialArray()...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range nums.Enumerate(nil).Select(sleepIdentity) {
			// Intentionally Left Blank
		}
	}
}

func sleepIdentity(num interface{}) interface{} {
	time.Sleep(2 * time.Millisecond)
	return Identity(num)
}

func getInitializedSequentialArray() []interface{} {

	rawNums := make([]interface{}, 1000, 1000)
	for i := range rawNums {
		rawNums[i] = i + 1
	}
	return rawNums
}
