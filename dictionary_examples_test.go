package collection_test

import (
	"context"
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
	fmt.Println(collection.CountAll[string](subject))

	subject.Clear()

	fmt.Println(subject.Size())
	fmt.Println(collection.Any[string](subject))

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

	upperCase := collection.Select[string](subject, strings.ToUpper)

	for word := range subject.Enumerate(context.Background()) {
		fmt.Println(word)
	}

	for word := range upperCase.Enumerate(context.Background()) {
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
	fmt.Println(collection.CountAll[string](subject))

	subject.Remove(world)

	fmt.Println(subject.Size())
	fmt.Println(collection.CountAll[string](subject))
	fmt.Println(collection.Any[string](subject))

	// Output:
	// 2
	// 2
	// 1
	// 1
	// true
}
