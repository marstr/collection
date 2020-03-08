package collection

type LRUCache struct {
	capacity uint
	entries map[interface{}]lruEntry
	touched LinkedList
}

type lruEntry struct {
	Priority *llNode
	Payload interface{}
}

func (lru *LRUCache) Put(key interface{}, value interface{}) bool {

}

func (lru *LRUCache) Get(key interface{}) (interface{}, bool) {
	if result, ok := lru.entries[key]; ok {
		lru.touched.first
	}
	return nil, false
}