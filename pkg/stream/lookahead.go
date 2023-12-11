package stream

import (
	"io"

	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/ringbuffer"
	"github.com/cockroachdb/errors"
)

// Lookahead returns a stream in which we can look ahead
// from the current point
func Lookahead[V any](input Stream[V]) *LookaheadStream[V] {
	return &LookaheadStream[V]{
		input:  input,
		buffer: ringbuffer.NewGrowable[V](),
	}
}

type LookaheadStream[V any] struct {
	input  Stream[V]
	buffer *ringbuffer.Growable[V]
}

func (l *LookaheadStream[V]) Next() (next V, err error) {
	if buffered, valid := l.buffer.Dequeue(); valid {
		return buffered, nil
	}

	return l.input.Next()
}

func (l *LookaheadStream[V]) PeekN(n int) (peek []V, err error) {
	for l.buffer.Len() < n {
		next, err := l.input.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return l.buffer.PeekN(n), nil
			} else {
				return nil, errors.Wrap(err, "unable to peak")
			}
		} else {
			l.buffer.Push(next)
		}
	}

	// Return the buffer
	return l.buffer.PeekN(n), nil
}
