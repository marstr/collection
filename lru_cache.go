package collection

import (
	"context"
	"sync"
)

// LRUCache hosts up to a given number of items. When more are presented, the least recently used item
// is evicted from the cache.
type LRUCache[K comparable, V any] struct {
	capacity uint
	entries  map[K]*lruEntry[K, V]
	touched  *LinkedList[*lruEntry[K, V]]
	key      sync.RWMutex
}

type lruEntry[K any, V any] struct {
	Node  *llNode[*lruEntry[K, V]]
	Key   K
	Value V
}

// NewLRUCache creates an empty cache, which will accommodate the given number of items.
func NewLRUCache[K comparable, V any](capacity uint) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		entries:  make(map[K]*lruEntry[K, V], capacity+1),
		touched:  NewLinkedList[*lruEntry[K, V]](),
	}
}

// Put adds a value to the cache. The added value may be expelled without warning.
func (lru *LRUCache[K, V]) Put(key K, value V) {
	lru.key.Lock()
	defer lru.key.Unlock()

	entry, ok := lru.entries[key]
	if ok {
		lru.touched.removeNode(entry.Node)
	} else {
		entry = &lruEntry[K, V]{
			Node: &llNode[*lruEntry[K, V]]{},
			Key:  key,
		}
	}

	entry.Node.payload = entry
	entry.Value = value
	lru.touched.addNodeFront(entry.Node)
	lru.entries[key] = entry

	if lru.touched.Length() > lru.capacity {
		removed, ok := lru.touched.RemoveBack()
		if ok {
			delete(lru.entries, removed.Key)
		}
	}
}

// Get retrieves a cached value, if it is still present.
func (lru *LRUCache[K, V]) Get(key K) (V, bool) {
	lru.key.RLock()
	defer lru.key.RUnlock()

	entry, ok := lru.entries[key]
	if !ok {
		return *new(V), false
	}

	lru.touched.removeNode(entry.Node)
	lru.touched.addNodeFront(entry.Node)
	return entry.Node.payload.Value, true
}

// Remove explicitly takes an item out of the cache.
func (lru *LRUCache[K, V]) Remove(key K) bool {
	lru.key.RLock()
	defer lru.key.RUnlock()

	entry, ok := lru.entries[key]
	if !ok {
		return false
	}

	lru.touched.removeNode(entry.Node)
	delete(lru.entries, key)
	return true
}

// Enumerate lists each value in the cache.
func (lru *LRUCache[K, V]) Enumerate(ctx context.Context) Enumerator[V] {
	retval := make(chan V)

	nested := lru.touched.Enumerate(ctx)

	go func() {
		lru.key.RLock()
		defer lru.key.RUnlock()
		defer close(retval)

		for entry := range nested {
			select {
			case retval <- entry.Value:
				break
			case <-ctx.Done():
				return
			}
		}
	}()

	return retval
}

// EnumerateKeys lists each key in the cache.
func (lru *LRUCache[K, V]) EnumerateKeys(ctx context.Context) Enumerator[K] {
	retval := make(chan K)

	nested := lru.touched.Enumerate(ctx)

	go func() {
		lru.key.RLock()
		defer lru.key.RUnlock()
		defer close(retval)

		for entry := range nested {
			select {
			case retval <- entry.Key:
				break
			case <-ctx.Done():
				return
			}
		}
	}()

	return retval
}
