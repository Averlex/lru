package lru

import "sync"

// Cache is an interface for an LRU cache.
type Cache[K comparable, V any] interface {
	Set(key K, value V) bool
	Get(key K) (V, bool)
	Clear()
}

type lruCache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	queue    List[*cacheListItem[K, V]]
	items    map[K]*ListItem[*cacheListItem[K, V]]
}

type cacheListItem[K comparable, V any] struct {
	key   K
	value V
}

// NewCache returns a new Cache with the given capacity. If the capacity is less than 1, it returns nil.
// The cache is implemented as a doubly-linked list with a map from keys to list items.
func NewCache[K comparable, V any](capacity int) Cache[K, V] {
	if capacity < 1 {
		return nil
	}

	return &lruCache[K, V]{
		capacity: capacity,
		queue:    NewList[*cacheListItem[K, V]](),
		items:    make(map[K]*ListItem[*cacheListItem[K, V]], capacity),
	}
}

// Set adds a key-value pair to the cache. If the key already exists, it updates the value
// and moves the item to the front of the queue. If the cache exceeds its capacity, it removes
// the least recently used item. Returns true if the key was already present in the cache, false otherwise.
func (c *lruCache[K, V]) Set(key K, value V) bool {
	listItem := &cacheListItem[K, V]{key, value}

	c.mu.Lock()
	defer c.mu.Unlock()

	// The element is present in the cache -> updating it's value, moving it to the front.
	if v, ok := c.items[key]; ok {
		v.Value = listItem
		c.queue.MoveToFront(v)
		return true
	}

	newElem := c.queue.PushFront(listItem)
	c.items[key] = newElem

	// Removing the oldest cache item to sustain the capacity.
	if c.queue.Len() > c.capacity {
		delete(c.items, c.queue.Back().Value.key)
		c.queue.Remove(c.queue.Back())
	}

	return false
}

// Get returns a value for a key if it exists in the cache, also moves the accessed item
// to the front of the queue. Otherwise, returns zero value and false.
func (c *lruCache[K, V]) Get(key K) (V, bool) {
	var zeroVal V

	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok := c.items[key]; ok {
		c.queue.MoveToFront(v)
		return v.Value.value, true
	}

	return zeroVal, false
}

// Clear removes all stored items from the cache.
func (c *lruCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList[*cacheListItem[K, V]]()
	c.items = make(map[K]*ListItem[*cacheListItem[K, V]], c.capacity)
}
