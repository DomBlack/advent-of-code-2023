package stream

import (
	"bufio"
	"bytes"
	"io"
)

// From returns a stream from a given slice
func From[V any](input []V) Stream[V] {
	return &sliceSource[V]{
		values:  input,
		nextIdx: 0,
	}
}

func FromItem[V any](input V) Stream[V] {
	return &sliceSource[V]{
		values:  []V{input},
		nextIdx: 0,
	}
}

type sliceSource[V any] struct {
	values  []V
	nextIdx int
}

func (l *sliceSource[V]) Next() (next V, err error) {
	if l.nextIdx >= len(l.values) {
		return next, io.EOF
	}

	next = l.values[l.nextIdx]
	l.nextIdx++
	return next, nil
}

// LinesFrom returns a stream of lines from the given input.
func LinesFrom(input []byte) Stream[string] {
	return lineSource{
		scanner: bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(input))),
	}
}

type lineSource struct {
	scanner *bufio.Scanner
}

func (l lineSource) Next() (string, error) {
	if !l.scanner.Scan() {
		err := l.scanner.Err()
		if err != nil {
			return "", err
		} else {
			return "", io.EOF
		}
	}

	return l.scanner.Text(), nil
}
