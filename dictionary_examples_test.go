package collection_test

import (
	"fmt"
	"strings"

	"github.com/marstr/collection/v2"
)

func ExampleDictionary_Add() {
	subject := &collection.Dictionary{}

	const example = "hello"
	fmt.Println(subject.Contains(example))
	fmt.Println(subject.Size())
	subject.Add(example)
	fmt.Println(subject.Contains(example))
	fmt.Println(subject.Size())

	// Output:
	// false
	// 0
	// true
	// 1
}

func ExampleDictionary_Clear() {
	subject := &collection.Dictionary{}

	subject.Add("hello")
	subject.Add("world")

	fmt.Println(subject.Size())
	fmt.Println(collection.CountAll(subject))

	subject.Clear()

	fmt.Println(subject.Size())
	fmt.Println(collection.Any(subject))

	// Output:
	// 2
	// 2
	// 0
	// false
}

func ExampleDictionary_Enumerate() {
	subject := &collection.Dictionary{}
	subject.Add("world")
	subject.Add("hello")

	upperCase := collection.Select(subject, func(x interface{}) interface{} {
		return strings.ToUpper(x.(string))
	})

	for word := range subject.Enumerate(nil) {
		fmt.Println(word)
	}

	for word := range upperCase.Enumerate(nil) {
		fmt.Println(word)
	}

	// Output:
	// hello
	// world
	// HELLO
	// WORLD
}

func ExampleDictionary_Remove() {
	const world = "world"
	subject := &collection.Dictionary{}
	subject.Add("hello")
	subject.Add(world)

	fmt.Println(subject.Size())
	fmt.Println(collection.CountAll(subject))

	subject.Remove(world)

	fmt.Println(subject.Size())
	fmt.Println(collection.CountAll(subject))
	fmt.Println(collection.Any(subject))

	// Output:
	// 2
	// 2
	// 1
	// 1
	// true
}
