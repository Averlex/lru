// Package lru provides a doubly-linked list and LRU cache implementations.
package lru

// List represents some basic operations over a doubly-linked list.
type List[V any] interface {
	Len() int
	Front() *ListItem[V]
	Back() *ListItem[V]
	PushFront(v V) *ListItem[V]
	PushBack(v V) *ListItem[V]
	Remove(elem *ListItem[V])
	MoveToFront(elem *ListItem[V])
}

// ListItem represents a basic item of the doubly-linked list.
type ListItem[V any] struct {
	Value V
	Next  *ListItem[V]
	Prev  *ListItem[V]
}

type list[V any] struct {
	len   int
	front *ListItem[V]
	back  *ListItem[V]
}

// NewList returns a new list with 0 length.
func NewList[V any]() List[V] {
	return new(list[V])
}

// Len returns the length of the list.
func (l *list[V]) Len() int {
	return l.len
}

// Front returns the first item in the list.
// If the list is empty, it returns nil.
func (l *list[V]) Front() *ListItem[V] {
	return l.front
}

// Back returns the last item in the list.
// If the list is empty, it returns nil.
func (l *list[V]) Back() *ListItem[V] {
	return l.back
}

// pushFrontLogic is a helper method for PushFront and MoveToFront methods.
func (l *list[V]) pushFrontLogic(i *ListItem[V]) *ListItem[V] {
	switch l.len {
	case 0:
		l.front = i
		l.back = i
		l.front.Next = nil
	case 1:
		l.front = i
		l.front.Next = l.back
		l.back.Prev = l.front
	default:
		l.front.Prev = i
		i.Next = l.front
		l.front = i
	}

	l.front.Prev = nil
	l.len++

	return i
}

// PushFront adds the value v at the beginning of the list.
// The function returns the item that was created for the value v.
func (l *list[V]) PushFront(v V) *ListItem[V] {
	newItem := &ListItem[V]{Value: v}

	return l.pushFrontLogic(newItem)
}

// PushBack adds the value v at the end of the list.
// The function returns the item that was created for the value v.
func (l *list[V]) PushBack(v V) *ListItem[V] {
	newItem := &ListItem[V]{Value: v}

	switch l.len {
	case 0:
		l.front = newItem
		l.back = newItem
	case 1:
		l.back = newItem
		l.back.Prev = l.front
		l.front.Next = l.back
	default:
		l.back.Next = newItem
		newItem.Prev = l.back
		l.back = newItem
	}

	l.len++

	return newItem
}

// Remove deletes the specified ListItem from the list.
// The length of the list is decremented by one.
func (l *list[V]) Remove(elem *ListItem[V]) {
	switch l.len {
	case 0:
		return
	case 1:
		// Assuming the list item is in the given list - panic otherwise.
		l.front = nil
		l.back = nil
	default:
		switch elem {
		case l.front:
			l.front = elem.Next
			l.front.Prev = nil
		case l.back:
			l.back = elem.Prev
			l.back.Next = nil
		// Item is somewhere in the middle of the list.
		default:
			// Assuming the list item is in the given list - panic otherwise.
			elem.Prev.Next = elem.Next
			elem.Next.Prev = elem.Prev
		}
	}

	l.len--
}

// MoveToFront moves item i to the front of the list.
// For an empty list function has the similar behavior as PushFront method.
func (l *list[V]) MoveToFront(elem *ListItem[V]) {
	l.Remove(elem)
	l.pushFrontLogic(elem)
}
