package collection

import (
	"context"
	"strings"
	"testing"
)

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
