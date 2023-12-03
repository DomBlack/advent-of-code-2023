package stream

func Resettable[V any](input Stream[V]) ResettableStream[V] {
	return &resettableStream[V]{input: input}
}

type ResettableStream[V any] interface {
	Stream[V]
	Save()    // Save the current position of the stream
	Restore() // Restore the stream to the last saved position (if no save has been made then this will be the beginning)
	Reset()   // Reset the stream to the beginning
}

type resettableStream[V any] struct {
	input     Stream[V]
	collected []V
	idx       int
	saveIdx   int
}

func (r *resettableStream[V]) Next() (V, error) {
	if r.idx < len(r.collected) {
		v := r.collected[r.idx]
		r.idx++
		return v, nil
	}

	v, err := r.input.Next()
	if err != nil {
		return v, err
	}

	r.collected = append(r.collected, v)
	r.idx++

	return v, nil
}

func (r *resettableStream[V]) Reset() {
	r.idx = 0
}

// Save saves the current position of the stream so that it can be restored later
func (r *resettableStream[V]) Save() {
	r.saveIdx = r.idx
}

// Restore restores the stream to the last saved position
func (r *resettableStream[V]) Restore() {
	r.idx = r.saveIdx
}
