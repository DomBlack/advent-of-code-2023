package stream

// Map returns a new stream with the given function applied to each element.
//
// If the function returns an error the stream will stop and the error will be
// returned.
func Map[A, B any](input Stream[A], fn func(A) (B, error)) Stream[B] {
	return mapStream[A, B]{
		input: input,
		fn:    fn,
	}
}

type mapStream[A, B any] struct {
	input Stream[A]
	fn    func(A) (B, error)
}

func (m mapStream[A, B]) Next() (next B, err error) {
	input, err := m.input.Next()
	if err != nil {
		return next, err
	}

	return m.fn(input)
}

// Filter returns a stream with only the elements that match the given predicate.
//
// If the predicate returns an error the stream will stop and the error will be
func Filter[A any](input Stream[A], predicate func(A) (keep bool, err error)) Stream[A] {
	return filterStream[A]{
		input: input,
		fn:    predicate,
	}
}

type filterStream[V any] struct {
	input Stream[V]
	fn    func(V) (bool, error)
}

func (f filterStream[V]) Next() (V, error) {
	for {
		input, err := f.input.Next()
		if err != nil {
			return input, err
		}

		keep, err := f.fn(input)
		if err != nil {
			return input, err
		}

		if keep {
			return input, nil
		}
	}
}
