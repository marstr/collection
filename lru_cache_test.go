package collection

import "testing"

func TestLRUCache_Put_replace(t *testing.T) {
	const key = 1
	const firstPut = "first"
	const secondPut = "second"

	subject := NewLRUCache[int, string](10)
	subject.Put(key, firstPut)
	subject.Put(key, secondPut)

	want := secondPut
	got, ok := subject.Get(key)
	if !ok {
		t.Logf("key should have been present")
		t.Fail()
	}

	if got != want {
		t.Logf("Unexpected result\n\tgot:  %s\n\twant: %s", got, want)
		t.Fail()
	}
}

func TestLRUCache_Remove_empty(t *testing.T) {
	subject := NewLRUCache[int, int](10)
	got := subject.Remove(7)
	if got != false {
		t.Fail()
	}
}

func TestLRUCache_Remove_present(t *testing.T) {
	const key = 10
	subject := NewLRUCache[int, string](6)
	subject.Put(key, "ten")
	ok := subject.Remove(key)
	if !ok {
		t.Fail()
	}

	_, ok = subject.Get(key)
	if ok {
		t.Fail()
	}
}

func TestLRUCache_Remove_notPresent(t *testing.T) {
	const key1 = 10
	const key2 = key1 + 1
	subject := NewLRUCache[int, string](6)
	subject.Put(key2, "eleven")
	ok := subject.Remove(key1)
	if ok {
		t.Fail()
	}

	_, ok = subject.Get(key2)
	if !ok {
		t.Fail()
	}
}
