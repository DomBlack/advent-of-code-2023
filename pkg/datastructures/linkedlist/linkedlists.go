package linkedlist

type LinkedList[T any] struct {
	first  *node[T]
	last   *node[T]
	length int
}

type node[T any] struct {
	value T
	next  *node[T]
	prev  *node[T]
}

// New creates a doubly linked list
func New[T any](items ...T) *LinkedList[T] {
	list := &LinkedList[T]{}
	for _, item := range items {
		list.Push(item)
	}
	return list
}

// Push adds a value to the end of the list
func (l *LinkedList[T]) Push(value T) {
	if l.first == nil {
		l.first = &node[T]{value: value}
		l.last = l.first
		l.length = 1
	} else {
		l.last.next = &node[T]{
			prev:  l.last,
			value: value,
		}
		l.last = l.last.next
		l.length += 1
	}
}

// Pop removes a value from the end of the list
func (l *LinkedList[T]) Pop() (value T, valid bool) {
	if l.last == nil {
		return value, false
	}

	value = l.last.value
	l.last = l.last.prev
	l.last.next = nil
	l.length -= 1

	// If we removed the last item, clear the first pointer
	if l.last == nil {
		l.first = nil
	}

	return value, true
}

// Dequeue returns the first value on the list
func (l *LinkedList[T]) Dequeue() (value T, valid bool) {
	if l.first == nil {
		return value, false
	}

	value = l.first.value
	l.first = l.first.next
	l.first.prev = nil
	l.length -= 1

	// If we removed the last item, clear the last pointer
	if l.first == nil {
		l.last = nil
	}

	return value, true
}
