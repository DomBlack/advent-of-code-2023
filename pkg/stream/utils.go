package stream

import (
	"github.com/rs/zerolog"
)

func DebugPrint[V any](log zerolog.Logger, stageName string, stream Stream[V]) Stream[V] {
	if !log.Debug().Enabled() {
		return stream
	}

	values, err := Collect(stream)
	if err != nil {
		return &errStream[V]{
			err: err,
		}
	}

	log.Debug().Msgf("%s values", stageName)
	for _, value := range values {
		log.Debug().Msgf("  %v", value)
	}

	return From(values)
}

type errStream[V any] struct {
	err error
}

func (e errStream[V]) Next() (next V, err error) {
	return next, e.err
}
