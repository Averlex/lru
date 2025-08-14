# lru

A simple and effective implementation of concurrent-safe LRU cache on Go.

## Features

- ✅ O(1) average time complexity for `Get` and `Set`
- ✅ Thread-safe with `sync.Mutex`
- ✅ Automatic eviction of least recently used items
- ✅ Clear operation to reset the cache
- ✅ Simple, idiomatic Go interface

## Usage

```go
import "github.com/yourname/lru"

cache := lru.NewCache(100) // capacity = 100

cache.Set("key1", "value1")
if val, ok := cache.Get("key1"); ok {
    fmt.Println("Found:", val)
}

cache.Clear()
```

## Interface

```go
type Cache interface {
    Set(key Key, value any) bool
    Get(key Key) (any, bool)
    Clear()
}
```

- `Set` returns `true` if the key already existed.
- `Get` returns the value and a boolean indicating presence.
- `Clear` removes all entries.

## Implementation

- Uses a `map[Key]*ListItem` for O(1) access.
- Doubly-linked list (`List`) to maintain access order.
- Guarded by a mutex for concurrent access.

## Installation

```bash
go get github.com/yourname/lru
```
