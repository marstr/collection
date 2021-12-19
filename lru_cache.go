package collection

import (
	"context"
	"sync"
)

// LRUCache hosts up to a given number of items. When more are presented, the least recently used item
// is evicted from the cache.
type LRUCache struct {
	capacity uint
	entries  map[interface{}]*lruEntry
	touched  *LinkedList
	key      sync.RWMutex
}

type lruEntry struct {
	Node  *llNode
	Key   interface{}
	Value interface{}
}

// NewLRUCache creates an empty cache, which will accommodate the given number of items.
func NewLRUCache(capacity uint) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		entries:  make(map[interface{}]*lruEntry, capacity+1),
		touched:  NewLinkedList(),
	}
}

// Put adds a value to the cache. The added value may be expelled without warning.
func (lru *LRUCache) Put(key interface{}, value interface{}) {
	lru.key.Lock()
	defer lru.key.Unlock()

	entry, ok := lru.entries[key]
	if ok {
		lru.touched.removeNode(entry.Node)
	} else {
		entry = &lruEntry{
			Node: &llNode{},
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
			delete(lru.entries, removed.(*lruEntry).Key)
		}
	}
}

// Get retrieves a cached value, if it is still present.
func (lru *LRUCache) Get(key interface{}) (interface{}, bool) {
	lru.key.RLock()
	defer lru.key.RUnlock()

	entry, ok := lru.entries[key]
	if !ok {
		return nil, false
	}

	lru.touched.removeNode(entry.Node)
	lru.touched.addNodeFront(entry.Node)
	return entry.Node.payload.(*lruEntry).Value, true
}

// Remove explicitly takes an item out of the cache.
func (lru *LRUCache) Remove(key interface{}) bool {
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
func (lru *LRUCache) Enumerate(ctx context.Context) Enumerator {
	retval := make(chan interface{})

	nested := lru.touched.Enumerate(ctx)

	go func() {
		lru.key.RLock()
		defer lru.key.RUnlock()
		defer close(retval)

		for entry := range nested {
			select {
			case retval <- entry.(*lruEntry).Value:
				break
			case <-ctx.Done():
				return
			}
		}
	}()

	return retval
}

// EnumerateKeys lists each key in the cache.
func (lru *LRUCache) EnumerateKeys(ctx context.Context) Enumerator {
	retval := make(chan interface{})

	nested := lru.touched.Enumerate(ctx)

	go func() {
		lru.key.RLock()
		defer lru.key.RUnlock()
		defer close(retval)

		for entry := range nested {
			select {
			case retval <- entry.(*lruEntry).Key:
				break
			case <-ctx.Done():
				return
			}
		}
	}()

	return retval
}
