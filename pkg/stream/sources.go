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
	return scannerSource{
		scanner: bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(input))),
	}
}

// SplitBy returns a stream of strings split by the given byte.
func SplitBy(input []byte, split byte) Stream[string] {
	s := scannerSource{
		scanner: bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(input))),
	}
	s.scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, split); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})
	return s
}

type scannerSource struct {
	scanner *bufio.Scanner
}

func (l scannerSource) Next() (string, error) {
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
