# collection
[![PkgGoDev](https://pkg.go.dev/badge/github.com/marstr/collection/v2)](https://pkg.go.dev/github.com/marstr/collection/v2) [![Build and Test](https://github.com/marstr/collection/workflows/Build%20and%20Test/badge.svg)](https://github.com/marstr/collection/actions?query=workflow%3A"Build+and+Test")

# Usage

## Available Data Structures:

### Dictionary
This is a logical set of strings. It utilizes a prefix tree model to be very space efficient.

### LinkedList
A collection that offers fast, consistent insertion time when adding to either the beginning or end. Accessing a random element is slower than other similar list data structures.

### List
Similar to a C++ `Vector`, Java `ArrayList`, or C# `List` this is a wrapper over top of arrays that allows for quick random access, but somewhat slower insertion characteristics than a `LinkedList`.

### LRUCache
This name is short for "Least Recently Used Cache". It holds a predetermined number of items, and as new items inserted, the least recently added or read item will be removed. This can be a useful way to build a tool that uses the proxy pattern to have quick access to the most useful items, and slower access to any other item. There is a memory cost for this, but it's often worth it.

### Queue
Stores items without promising random access. The first thing you put in will be the first thing you get out.

### Stack
Stores items without promising random access. The first thing you put in will be the last thing you get out.

## Querying Collections
Inspired by .NET's Linq, querying data structures used in this library is a snap! Build Go pipelines quickly and easily which will apply lambdas as they query your data structures.

### Slices
Converting between slices and a queryable structure is as trivial as it should be.
``` Go
original := []string{"a", "b", "c"}
subject := collection.AsEnumerable(original...)

for entry := range subject.Enumerate(context.Background()) {
    fmt.Println(entry)
}
// Output:
// a
// b
// c

```

### Where
``` Go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

subject := collection.AsEnumerable[int](1, 2, 3, 4, 5, 6)
filtered := collection.Where(subject, func(num int) bool{
    return num > 3
})
for entry := range filtered.Enumerate(ctx) {
    fmt.Println(entry)
}
// Output:
// 4
// 5
// 6
```
### Select
``` Go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

subject := collection.AsEnumerable[int](1, 2, 3, 4, 5, 6)
updated := collection.Select[int](subject, func(num int) int {
    return num + 10
})
for entry := range updated.Enumerate(ctx) {
    fmt.Println(entry)
}

// Output:
// 11
// 12
// 13
// 14
// 15
// 16
```

## Queues
### Creating a Queue

``` Go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

subject := collection.NewQueue(1, 2, 3, 5, 8, 13, 21)
selected := subject.Enumerate(ctx).Skip(3).Take(3)
for entry := range selected {
	fmt.Println(entry)
}

// Output:
// 5
// 8
// 13
```

### Checking if a Queue is empty
``` Go
populated := collection.NewQueue(1, 2, 3, 5, 8, 13)
notPopulated := collection.NewQueue[int]()
fmt.Println(populated.IsEmpty())
fmt.Println(notPopulated.IsEmpty())
// Output:
// false
// true
```

## Other utilities

### Fibonacci
This was added to test Enumerable types that have no logical conclusion. But it may prove useful other places, so it is available in the user-facing package and not hidden away in a test package.

### Filesystem
Find the standard library's pattern for looking through a directory cumbersome? Use the collection querying mechanisms seen above to search a directory as a collection of files and child directories.

# Versioning
This library will conform to strict semantic versions as defined by [semver.org](http://semver.org/spec/v2.0.0.html)'s v2 specification.

# Contributing
I accept contributions! Please submit PRs to the `main` or `v1` branches. Remember to add tests!

# F.A.Q.

## Should I use v1 or v2?

If you are newly adopting this library, and are able to use Go 1.18 or newer, it is highly recommended that you use v2.

V2 was primarily added to support Go generics when they were introduced in Go 1.18, but there were other breaking changes made because of the opportunity to do with the major version bump.

Because it's not reasonable to expect everybody to adopt the newest versions of Go immediately as they're released, v1 of this library wil be activey supported until Go 1.17 is no longer supported by the Go team. After that community contributions to v1 will be entertained, but active development won't be ported to the `v1` branch.

## Why does `Enumerate` take a `context.Context`?

Having a context associated with the enumeration allows for cancellation. This is valuable in some scenarios, where enumeration may be a time-consuming operation. For example, imagine an `Enumerable` that wraps a web API which returns results in pages. Injecting a context
allows for you to add operation timeouts, and otherwise protect yourself from an operation that may not finish quickly enough for you (or at all.)

However, under the covers an Enumerator[T] is a `<-chan T`. This decision means that a separate goroutine is used to publish to the channel while your goroutine reads from it.

**That means if your code stops before all items in the Enumerator are read, a goroutine and all of the memory it's using will be leaked.**

This is a known problem, and it's understood why it's not ideal. The workaround is easy - if there's ever a chance you won't enumerate all items, protect yourself by using the following pattern:

``` Go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// ...

for item := range myEnumerable.Enumerate(ctx) {
    // ...
}
```