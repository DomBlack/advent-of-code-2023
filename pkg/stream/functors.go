package stream

import (
	"io"

	"github.com/cockroachdb/errors"
)

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

// FlatMap returns a new merged stream with the given function applied to each element
// and the results merged into a single stream.
//
// If the function returns an error the stream will stop and the error will be
// returned.
func FlatMap[A, B any](input Stream[A], fn func(A) (Stream[B], error)) Stream[B] {
	return &flatMapStream[A, B]{
		input: input,
		fn:    fn,
	}
}

type flatMapStream[A, B any] struct {
	input Stream[A]
	temp  Stream[B]
	fn    func(A) (Stream[B], error)
}

func (m *flatMapStream[A, B]) Next() (next B, err error) {
	// If we have a temp stream then we need to get the next element from that
	if m.temp != nil {
		next, err := m.temp.Next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return next, errors.WithStack(err)
			} else {
				m.temp = nil
			}
		} else {
			return next, nil
		}
	}

	// Otherwise get the next element from the input stream
	input, err := m.input.Next()
	if err != nil {
		return next, err
	}

	// Map it into a new output stream
	m.temp, err = m.fn(input)
	if err != nil {
		return next, errors.WithStack(err)
	}

	// And return the first element from that
	return m.Next()
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
