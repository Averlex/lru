// Package lru provides a doubly-linked list and LRU cache implementations.
package lru

// List represents some basic operations over a doubly-linked list.
type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// ListItem represents a basic item of the doubly-linked list.
type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

// NewList returns a new list with 0 length.
func NewList() List {
	return new(list)
}

// Len returns the length of the list.
func (l *list) Len() int {
	return l.len
}

// Front returns the first item in the list.
// If the list is empty, it returns nil.
func (l *list) Front() *ListItem {
	return l.front
}

// Back returns the last item in the list.
// If the list is empty, it returns nil.
func (l *list) Back() *ListItem {
	return l.back
}

// pushFrontLogic is a helper method for PushFront and MoveToFront methods.
func (l *list) pushFrontLogic(i *ListItem) *ListItem {
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
func (l *list) PushFront(v any) *ListItem {
	newItem := &ListItem{Value: v}

	return l.pushFrontLogic(newItem)
}

// PushBack adds the value v at the end of the list.
// The function returns the item that was created for the value v.
func (l *list) PushBack(v any) *ListItem {
	newItem := &ListItem{Value: v}

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
func (l *list) Remove(i *ListItem) {
	switch l.len {
	case 0:
		return
	case 1:
		// Assuming the list item is in the given list - panic otherwise.
		l.front = nil
		l.back = nil
	default:
		switch i {
		case l.front:
			l.front = i.Next
			l.front.Prev = nil
		case l.back:
			l.back = i.Prev
			l.back.Next = nil
		// Item is somewhere in the middle of the list.
		default:
			// Assuming the list item is in the given list - panic otherwise.
			i.Prev.Next = i.Next
			i.Next.Prev = i.Prev
		}
	}

	l.len--
}

// MoveToFront moves item i to the front of the list.
// For an empty list function has the similar behavior as PushFront method.
func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.pushFrontLogic(i)
}
