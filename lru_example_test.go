package collection_test

import (
	"context"
	"fmt"
	"github.com/marstr/collection"
)

func ExampleLRUCache() {
	subject := collection.NewLRUCache[int, string](3)
	subject.Put(1, "one")
	subject.Put(2, "two")
	subject.Put(3, "three")
	subject.Put(4, "four")
	fmt.Println(subject.Get(1))
	fmt.Println(subject.Get(4))
	// Output:
	// false
	// four true
}

func ExampleLRUCache_Enumerate() {
	ctx := context.Background()
	subject := collection.NewLRUCache[int, string](3)
	subject.Put(1, "one")
	subject.Put(2, "two")
	subject.Put(3, "three")
	subject.Put(4, "four")

	for key := range subject.Enumerate(ctx.Done()) {
		fmt.Println(key)
	}

	// Output:
	// four
	// three
	// two
}

func ExampleLRUCache_EnumerateKeys() {
	ctx := context.Background()
	subject := collection.NewLRUCache[int, string](3)
	subject.Put(1, "one")
	subject.Put(2, "two")
	subject.Put(3, "three")
	subject.Put(4, "four")

	for key := range subject.EnumerateKeys(ctx.Done()) {
		fmt.Println(key)
	}

	// Output:
	// 4
	// 3
	// 2
}
