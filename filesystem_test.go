package collection

import (
	"context"
	"fmt"
	"math"
	"path"
	"path/filepath"
	"testing"
)

func TestEnumerateDirectoryOptions_UniqueBits(t *testing.T) {
	isPowerOfTwo := func(subject float64) bool {
		a := math.Abs(math.Log2(subject))
		b := math.Floor(a)

		return a-b < .0000001
	}

	if !isPowerOfTwo(64) {
		t.Log("isPowerOfTwo decided 64 is not a power of two.")
		t.FailNow()
	}

	if isPowerOfTwo(91) {
		t.Log("isPowerOfTwo decided 91 is a power of two.")
		t.FailNow()
	}

	seen := make(map[DirectoryOptions]struct{})

	declared := []DirectoryOptions{
		DirectoryOptionsExcludeFiles,
		DirectoryOptionsExcludeDirectories,
		DirectoryOptionsRecursive,
	}

	for _, option := range declared {
		if _, ok := seen[option]; ok {
			t.Logf("Option: %d has already been declared.", option)
			t.Fail()
		}
		seen[option] = struct{}{}

		if !isPowerOfTwo(float64(option)) {
			t.Logf("Option should have been a power of two, got %g instead.", float64(option))
			t.Fail()
		}
	}
}

func ExampleDirectory_Enumerate() {
	traverser := Directory{
		Location: ".",
		Options:  DirectoryOptionsExcludeDirectories,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fileNames := Select[string](traverser, func(subject string) string {
		return path.Base(subject)
	})

	filesOfInterest := Where(fileNames, func(subject string) bool {
		return subject == "filesystem_test.go"
	})

	for entry := range filesOfInterest.Enumerate(ctx) {
		fmt.Println(entry)
	}

	// Output: filesystem_test.go
}

func TestDirectory_Enumerate(t *testing.T) {
	subject := Directory{
		Location: filepath.Join(".", "testdata", "foo"),
	}

	testCases := []struct {
		options  DirectoryOptions
		expected map[string]struct{}
	}{
		{
			options: 0,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "a.txt"): {},
				filepath.Join("testdata", "foo", "c.txt"): {},
				filepath.Join("testdata", "foo", "bar"):   {},
			},
		},
		{
			options: DirectoryOptionsExcludeFiles,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "bar"): {},
			},
		},
		{
			options: DirectoryOptionsExcludeDirectories,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "a.txt"): {},
				filepath.Join("testdata", "foo", "c.txt"): {},
			},
		},
		{
			options: DirectoryOptionsRecursive,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "bar"):          {},
				filepath.Join("testdata", "foo", "bar", "b.txt"): {},
				filepath.Join("testdata", "foo", "a.txt"):        {},
				filepath.Join("testdata", "foo", "c.txt"):        {},
			},
		},
		{
			options: DirectoryOptionsExcludeFiles | DirectoryOptionsRecursive,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "bar"): {},
			},
		},
		{
			options: DirectoryOptionsRecursive | DirectoryOptionsExcludeDirectories,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "a.txt"):        {},
				filepath.Join("testdata", "foo", "bar", "b.txt"): {},
				filepath.Join("testdata", "foo", "c.txt"):        {},
			},
		},
		{
			options:  DirectoryOptionsExcludeDirectories | DirectoryOptionsExcludeFiles,
			expected: map[string]struct{}{},
		},
		{
			options:  DirectoryOptionsExcludeFiles | DirectoryOptionsRecursive | DirectoryOptionsExcludeDirectories,
			expected: map[string]struct{}{},
		},
	}

	for _, tc := range testCases {
		subject.Options = tc.options
		t.Run(fmt.Sprintf("%d", uint(tc.options)), func(t *testing.T) {
			for entry := range subject.Enumerate(context.Background()) {
				if _, ok := tc.expected[entry]; !ok {
					t.Logf("unexpected result: %q", entry)
					t.Fail()
				}
				delete(tc.expected, entry)
			}

			if len(tc.expected) != 0 {
				for unseenFile := range tc.expected {
					t.Logf("missing file: %q", unseenFile)
				}
				t.Fail()
			}
		})
	}
}
