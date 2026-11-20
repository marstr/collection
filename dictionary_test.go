package collection

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewDictionaryFromEnumerable(t *testing.T) {
	testCases := []struct {
		name    string
		entries []string
		want    []string
	}{
		{name: "empty"},
		{name: "duplicates and prefixes", entries: []string{"beta", "alpha", "alphabet", "alpha"}, want: []string{"alpha", "alphabet", "beta"}},
		{name: "whole entries", entries: []string{"", "two words", " padded "}, want: []string{"", "two words", " padded "}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dict, err := NewDictionaryFromEnumerable(context.Background(), EnumerableSlice[string](tc.entries))
			if err != nil {
				t.Fatal(err)
			}
			assertDictionaryEntries(t, dict, tc.want)
		})
	}
}

func TestNewDictionaryFromReader(t *testing.T) {
	longLine := strings.Repeat("a", 70*1024)
	testCases := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "empty"},
		{name: "duplicates", input: "beta\nalpha\nalpha\n", want: []string{"alpha", "beta"}},
		{name: "no final newline", input: "alpha\nbeta", want: []string{"alpha", "beta"}},
		{name: "blank lines", input: "\nalpha\n\n", want: []string{"", "alpha"}},
		{name: "whitespace", input: "two words\n padded \na\tb\n", want: []string{"two words", " padded ", "a\tb"}},
		{name: "CRLF", input: "alpha\r\n\r\nbeta\r\n", want: []string{"alpha", "", "beta"}},
		{name: "bare carriage return", input: "alpha\r", want: []string{"alpha\r"}},
		{name: "long line", input: longLine + "\n", want: []string{longLine}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dict, err := NewDictionaryFromReader(context.Background(), strings.NewReader(tc.input))
			if err != nil {
				t.Fatal(err)
			}
			assertDictionaryEntries(t, dict, tc.want)
		})
	}
}

func TestNewDictionaryFromReader_DataWithEOF(t *testing.T) {
	reader := dictionaryReaderFunc(func(p []byte) (int, error) {
		return copy(p, "alpha\nbeta"), io.EOF
	})
	dict, err := NewDictionaryFromReader(context.Background(), reader)
	if err != nil {
		t.Fatal(err)
	}
	assertDictionaryEntries(t, dict, []string{"alpha", "beta"})
}

func TestNewDictionaryFromReader_Error(t *testing.T) {
	wantErr := errors.New("read failed")
	reader := dictionaryReaderFunc(func(p []byte) (int, error) {
		return copy(p, "alpha\nbeta"), wantErr
	})
	dict, err := NewDictionaryFromReader(context.Background(), reader)
	if dict != nil || !errors.Is(err, wantErr) {
		t.Fatalf("got (%v, %v), want (nil, %v)", dict, err, wantErr)
	}
}

func TestNewDictionaryFromFilename(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "dictionary.txt")
	if err := os.WriteFile(filename, []byte("beta\r\nalpha\r\nalpha"), 0600); err != nil {
		t.Fatal(err)
	}
	dict, err := NewDictionaryFromFilename(context.Background(), filename)
	if err != nil {
		t.Fatal(err)
	}
	assertDictionaryEntries(t, dict, []string{"alpha", "beta"})
	if err := os.Remove(filename); err != nil {
		t.Fatalf("removing input file after construction: %v", err)
	}

	dict, err = NewDictionaryFromFilename(context.Background(), filename)
	if dict != nil || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("got (%v, %v), want (nil, file-not-found error)", dict, err)
	}
}

func TestNewDictionary_ContextCanceled(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "dictionary.txt")
	if err := os.WriteFile(filename, []byte("alpha\n"), 0600); err != nil {
		t.Fatal(err)
	}
	constructors := []struct {
		name string
		new  func(context.Context) (*Dictionary, error)
	}{
		{"enumerable", func(ctx context.Context) (*Dictionary, error) {
			return NewDictionaryFromEnumerable(ctx, Empty[string]())
		}},
		{"reader", func(ctx context.Context) (*Dictionary, error) {
			return NewDictionaryFromReader(ctx, strings.NewReader("alpha\n"))
		}},
		{"filename", func(ctx context.Context) (*Dictionary, error) {
			return NewDictionaryFromFilename(ctx, filename)
		}},
	}
	for _, tc := range constructors {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			dict, err := tc.new(ctx)
			if dict != nil || !errors.Is(err, context.Canceled) {
				t.Fatalf("got (%v, %v), want (nil, context.Canceled)", dict, err)
			}
		})
	}
}

func TestNewDictionaryFromEnumerable_CanceledDuringEnumeration(t *testing.T) {
	for _, closeSource := range []bool{true, false} {
		name := "open source"
		if closeSource {
			name = "closed source"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			source := dictionaryEnumerableFunc(func(context.Context) Enumerator[string] {
				entries := make(chan string, 1)
				entries <- "alpha"
				if closeSource {
					close(entries)
				}
				cancel()
				return entries
			})
			dict, err := NewDictionaryFromEnumerable(ctx, source)
			if dict != nil || !errors.Is(err, context.Canceled) {
				t.Fatalf("got (%v, %v), want (nil, context.Canceled)", dict, err)
			}
		})
	}
}

func TestNewDictionaryFromReader_CanceledDuringRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reader := dictionaryReaderFunc(func(p []byte) (int, error) {
		cancel()
		return copy(p, "alpha"), io.EOF
	})
	dict, err := NewDictionaryFromReader(ctx, reader)
	if dict != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("got (%v, %v), want (nil, context.Canceled)", dict, err)
	}
}

func assertDictionaryEntries(t *testing.T, dict *Dictionary, want []string) {
	t.Helper()
	if dict == nil {
		t.Fatal("expected an allocated dictionary")
	}
	if got := dict.Size(); got != int64(len(want)) {
		t.Errorf("Size() = %d, want %d", got, len(want))
	}
	for _, entry := range want {
		if !dict.Contains(entry) {
			t.Errorf("missing entry %q", entry)
		}
	}
}

type dictionaryReaderFunc func([]byte) (int, error)

func (read dictionaryReaderFunc) Read(p []byte) (int, error) {
	return read(p)
}

type dictionaryEnumerableFunc func(context.Context) Enumerator[string]

func (enumerate dictionaryEnumerableFunc) Enumerate(ctx context.Context) Enumerator[string] {
	return enumerate(ctx)
}

func TestDictionary_Enumerate(t *testing.T) {
	dictSets := [][]string{
		{"alpha", "beta", "charlie"},
		{"also", "always"},
		{"canned", "beans"},
		{"duplicated", "duplicated", "after"},
	}

	for _, ds := range dictSets {
		t.Run("", func(t *testing.T) {
			subject := Dictionary{}
			expected := make(map[string]bool)
			added := 0
			for _, entry := range ds {
				if subject.Add(entry) {
					added++
				}
				expected[entry] = false
			}

			expectedSize := len(expected)

			if added != expectedSize {
				t.Logf("`Add` returned true %d times, expected %d times", added, expectedSize)
				t.Fail()
			}

			if subjectSize := CountAll[string](subject); subjectSize != expectedSize {
				t.Logf("`CountAll` returned %d elements, expected %d", subjectSize, expectedSize)
				t.Fail()
			}

			prev := ""
			for result := range subject.Enumerate(context.Background()) {
				t.Log(result)
				if alreadySeen, ok := expected[result]; !ok {
					t.Log("An unadded value was returned")
					t.Fail()
				} else if alreadySeen {
					t.Logf("\"%s\" was duplicated", result)
					t.Fail()
				}

				if stringle(result, prev) {
					t.Logf("Results \"%s\" and \"%s\" were not alphabetized.", prev, result)
					t.Fail()
				}
				prev = result

				expected[result] = true
			}
		})
	}
}

func TestDictionary_Add(t *testing.T) {
	subject := Dictionary{}

	subject.Add("word")

	if rootChildrenCount := len(subject.root.Children); rootChildrenCount != 1 {
		t.Logf("The root should only have one child, got %d instead.", rootChildrenCount)
		t.Fail()
	}

	if retreived, ok := subject.root.Children['w']; ok {
		leaf := retreived.Navigate("ord")
		if leaf == nil {
			t.Log("Unable to navigate from `w`")
			t.Fail()
		} else if !leaf.IsWord {
			t.Log("leaf should have been a word")
			t.Fail()
		}
	} else {
		t.Log("Root doesn't have child for `w`")
		t.Fail()
	}
}

func TestTrieNode_Navigate(t *testing.T) {
	leaf := trieNode{
		IsWord: true,
	}
	subject := trieNode{
		Children: map[rune]*trieNode{
			'a': {
				Children: map[rune]*trieNode{
					'b': {
						Children: map[rune]*trieNode{
							'c': &leaf,
						},
					},
				},
			},
		},
	}

	testCases := []struct {
		address  string
		expected *trieNode
	}{
		{"abc", &leaf},
		{"abd", nil},
		{"", &subject},
		{"a", subject.Children['a']},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			if result := subject.Navigate(tc.address); result != tc.expected {
				t.Logf("got: %v want: %v", result, tc.expected)
				t.Fail()
			}
		})
	}
}

func Test_stringle(t *testing.T) {
	testCases := []struct {
		left     string
		right    string
		expected bool
	}{
		{"a", "b", true},
		{"b", "a", false},
		{"a", "a", true},
		{"alpha", "b", true},
		{"a", "beta", true},
		{"alpha", "alpha", true},
		{"alpha", "alphabet", true},
		{"alphabet", "alpha", false},
		{"", "a", true},
		{"", "", true},
	}

	for _, tc := range testCases {
		t.Run(strings.Join([]string{tc.left, tc.right}, ","), func(t *testing.T) {
			if got := stringle(tc.left, tc.right); got != tc.expected {
				t.Logf("got: %v want: %v", got, tc.expected)
				t.Fail()
			}
		})
	}
}

func stringle(left, right string) bool {
	other := []byte(right)
	for i, letter := range []byte(left) {
		if i >= len(other) {
			return false
		}

		if letter > other[i] {
			return false
		} else if letter < other[i] {
			break
		}
	}
	return true
}
