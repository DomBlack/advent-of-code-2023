package stream

type Stream[V any] interface {
	// Next returns the next value in the stream
	// or if the stream has finished returns [io.EOF]
	Next() (V, error)
}
