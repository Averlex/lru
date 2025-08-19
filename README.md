# lru

[![Go version](https://img.shields.io/badge/go-1.24.2+-blue.svg)](https://golang.org)
[![Go Reference](https://pkg.go.dev/badge/github.com/Averlex/lru.svg)](https://pkg.go.dev/github.com/Averlex/lru)
[![Go Report Card](https://goreportcard.com/badge/github.com/Averlex/lru)](https://goreportcard.com/report/github.com/Averlex/lru)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A simple and effective implementation of thread-safe LRU cache on Go.

## Features

- ✅ O(1) average time complexity for `Get`, `Set` and `Clear`
- ✅ Thread-safe with `sync.Mutex`
- ✅ Automatic eviction of least recently used items
- ✅ **Generic types** for keys and values — works with any comparable key type
- ✅ Simple, idiomatic Go interface

## Usage

**Any comparable key type is supported**

```go
import "github.com/Averlex/lru"

cache := lru.NewCache[string, string](100) // capacity = 100

cache.Set("key1", "value1")
if val, ok := cache.Get("key1"); ok {
    fmt.Println("Found:", val)
}

cache.Clear()
```

## Interface

```go
type Cache[K comparable, V any] interface {
    Set(key K, value V) bool
    Get(key K) (V, bool)
    Clear()
}
```

- `Set` returns `true` if the key already exists.
- `Get` returns the value and a boolean indicating it's presence in the cache.
- `Clear` removes all entries from the cache.

## Implementation

- Uses a `map[key]*ListItem` for O(1) access.
- Doubly-linked list (`List`) to maintain access order.
- Generic `List` and `ListItem` types ensure type safety without `interface{}` assertions.
- Guarded by a mutex for concurrent access.

## Installation

```bash
go get github.com/Averlex/lru
```
