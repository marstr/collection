package collection_test

import (
	"context"
	"fmt"
	"github.com/marstr/collection/v2"
	"path"
)

func ExampleDirectory_Enumerate() {
	traverser := collection.Directory{
		Location: ".",
		Options:  collection.DirectoryOptionsExcludeDirectories,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fileNames := collection.Select[string](traverser, path.Base)

	filesOfInterest := collection.Where(fileNames, func(subject string) bool {
		return subject == "filesystem_examples_test.go"
	})

	for entry := range filesOfInterest.Enumerate(ctx) {
		fmt.Println(entry)
	}

	// Output: filesystem_examples_test.go
}
