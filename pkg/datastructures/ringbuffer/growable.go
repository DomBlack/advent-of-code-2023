package ringbuffer

// Growable represents a growable
// ring buffer in which we can add items to
// and it will resize to fit everything in it
//
// However as we pop and dequeue items from it
// it will never shrink in size. This is useful
// for storing temporary lists without needing
// to allocate new slices constantly
type Growable[T any] struct {
	buffer     []T // The underlying buffer
	peekBuffer []T // An additional buffer we use for peeking if needing
	headIdx    int // The index of the head (first element) within the buffer
	length     int // The current length of the buffer
}

// NewGrowable returns a growable ring buffer, that can be used as temporary storage
func NewGrowable[T any]() *Growable[T] {
	return &Growable[T]{
		buffer: make([]T, 8),
	}
}

func (rb *Growable[T]) grow() {
	newBuffer := make([]T, len(rb.buffer)*2)

	// Copy over the values, starting at position 0,
	// this means the free space will always be at the end
	for i := 0; i < rb.length; i++ {
		pos := (rb.headIdx + i) % len(rb.buffer)
		newBuffer[i] = rb.buffer[pos]
	}
	rb.headIdx = 0
	rb.buffer = newBuffer
}

// Len returns the current length of the ring buffer
func (rb *Growable[T]) Len() int {
	return rb.length
}

// Push adds the value to the end of the buffer
func (rb *Growable[T]) Push(value T) {
	// Check if we need to grow the ring buffer
	if rb.length >= len(rb.buffer) {
		rb.grow()
	}

	pos := (rb.headIdx + rb.length) % len(rb.buffer)
	rb.buffer[pos] = value
	rb.length++
}

// Dequeue removes the first value from the buffer
func (rb *Growable[T]) Dequeue() (value T, valid bool) {
	if rb.length <= 0 {
		return value, false
	}

	value = rb.buffer[rb.headIdx]

	if rb.length == 1 {
		// As a special case always reset the first index to 0
		rb.headIdx = 0
		rb.length = 0
	} else {
		rb.headIdx = (rb.headIdx + 1) % len(rb.buffer)
		rb.length--
	}

	return value, true
}

// PeekN returns upto N items from the buffer
//
// If n is larger than the length of the buffer, the full buffer
// will be returned.
//
// The slice it returns is not safe to be reused or modified
// as it may be taken from the internals of the ring buffer.
func (rb *Growable[T]) PeekN(n int) []T {
	if rb.length == 0 {
		return nil
	}

	n = min(n, rb.length)

	if rb.headIdx+n < len(rb.buffer) {
		// If the request doesn't wrap around in the buffer,
		// return that part of the buffer slice
		return rb.buffer[rb.headIdx : rb.headIdx+n]
	} else {
		// We need to use our peek buffer
		if cap(rb.peekBuffer) < n {
			rb.peekBuffer = make([]T, 0, n)
		}

		// Reset the peek buffer
		rb.peekBuffer = rb.peekBuffer[:0]

		for i := 0; i < n; i++ {
			pos := (rb.headIdx + i) % len(rb.buffer)
			rb.peekBuffer = append(rb.peekBuffer, rb.buffer[pos])
		}

		return rb.peekBuffer
	}
}
