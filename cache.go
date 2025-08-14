package lru

import "sync"

// Key is a wrapper type for string keys.
type Key string

// Cache is an interface for an LRU cache.
type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheListItem struct {
	key   Key
	value any
}

// NewCache returns a new Cache with the given capacity. If the capacity is less than 1, it returns nil.
// The cache is implemented as a doubly-linked list with a map from keys to list items.
func NewCache(capacity int) Cache {
	if capacity < 1 {
		return nil
	}

	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Set adds a key-value pair to the cache. If the key already exists, it updates the value
// and moves the item to the front of the queue. If the cache exceeds its capacity, it removes
// the least recently used item. Returns true if the key was already present in the cache, false otherwise.
func (c *lruCache) Set(key Key, value any) bool {
	listItem := &cacheListItem{key, value}

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
		delete(c.items, c.queue.Back().Value.(*cacheListItem).key)
		c.queue.Remove(c.queue.Back())
	}

	return false
}

// Get returns a value for a key if it exists in the cache, also moves the accessed item
// to the front of the queue. Otherwise, returns nil and false.
func (c *lruCache) Get(key Key) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok := c.items[key]; ok {
		c.queue.MoveToFront(v)
		return v.Value.(*cacheListItem).value, true
	}

	return nil, false
}

// Clear removes all stored items from the cache.
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
