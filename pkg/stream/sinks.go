package stream

import (
	"errors"
	"fmt"
	"io"
)

// numeric is a constraint that limits what type of streams can be used with certain sinks
// such as [Sum] or [SumToString].
type numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

// ForEach applies the given function to each value in the stream.
func ForEach[V any](input Stream[V], fn func(V) error) error {
	for {
		v, err := input.Next()

		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}

		err = fn(v)
		if err != nil {
			return err
		}
	}
}

// Collect returns all values from the stream as a slice.
func Collect[V any](input Stream[V]) ([]V, error) {
	var values []V

	for {
		v, err := input.Next()

		if errors.Is(err, io.EOF) {
			return values, nil
		} else if err != nil {
			return values, err
		}

		values = append(values, v)
	}
}

// Reduce returns a single value from the stream by applying the reducer function to each value in the stream.
func Reduce[V, A any](input Stream[V], reducer func(acc A, value V) (A, error)) (A, error) {
	var value A

	for {
		v, err := input.Next()

		if errors.Is(err, io.EOF) {
			return value, nil
		} else if err != nil {
			return value, err
		}

		value, err = reducer(value, v)
		if err != nil {
			return value, err
		}
	}
}

// Sum returns the sum of all values in the stream.
func Sum[V numeric](input Stream[V]) (V, error) {
	return Reduce[V, V](input, func(acc V, value V) (V, error) {
		return acc + value, nil
	})
}

// SumToString returns the sum of all values in the stream as a string.
func SumToString[V numeric](input Stream[V]) (string, error) {
	sum, err := Sum[V](input)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", sum), nil
}
