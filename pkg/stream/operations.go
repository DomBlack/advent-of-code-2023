package stream

import (
	"github.com/cockroachdb/errors"
)

// Partition splits a stream into two streams based on the result of the goLeft function
func Partition[V any](input Stream[V], goLeft func(V) (bool, error)) (left Stream[V], right Stream[V]) {
	leftBuffer := &partitionStream[V]{}
	rightBuffer := &partitionStream[V]{}

	pull := func() error {
		next, err := input.Next()
		if err != nil {
			return errors.WithStack(err)
		}

		left, err := goLeft(next)
		if err != nil {
			return errors.Wrap(err, "failed to partition stream")
		}

		if left {
			leftBuffer.buffer = append(leftBuffer.buffer, next)
		} else {
			rightBuffer.buffer = append(rightBuffer.buffer, next)
		}
		return nil
	}

	leftBuffer.pull = pull
	rightBuffer.pull = pull

	return leftBuffer, rightBuffer
}

type partitionStream[V any] struct {
	err    error
	pull   func() error
	buffer []V
}

func (b *partitionStream[V]) Next() (next V, err error) {
	if len(b.buffer) > 0 {
		next = b.buffer[0]
		b.buffer = b.buffer[1:]
		return next, nil
	}

	if b.err != nil {
		return next, b.err
	}

	for len(b.buffer) == 0 {
		b.err = b.pull()
		if b.err != nil {
			return next, b.err
		}
	}

	next = b.buffer[0]
	b.buffer = b.buffer[:0]

	return next, nil
}
