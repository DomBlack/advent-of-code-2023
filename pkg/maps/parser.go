package maps

import (
	"io"

	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
)

type streamParser[TileType Tile] struct {
	input     stream.Stream[string]
	options   []MapOption
	parseTile func(r rune) (TileType, error)

	lines [][]TileType
}

// NewParseFunc creates a new parser that will parse exactly one map from the given byte slice.
//
// If more than one map is found, an error will be returned. If you want to parse multiple maps
// then use [NewStreamingParseFunc].
func NewParseFunc[TileType Tile](parseTile func(r rune) (TileType, error), options ...MapOption) func([]byte) (*Map[TileType], error) {
	return func(data []byte) (*Map[TileType], error) {
		streamingParser := NewStreamingParseFunc(parseTile, options...)(data)
		maps, err := stream.Collect(streamingParser)
		if err != nil {
			return nil, err
		}

		if len(maps) != 1 {
			return nil, errors.Newf("expected 1 map, got %d", len(maps))
		}

		return maps[0], nil
	}
}

// NewStreamingParseFunc creates a new stream parser that will parse one or more maps from an original
// byte slice.
//
// If you only need to parse a single map use [NewParseFunc].
func NewStreamingParseFunc[TileType Tile](parseTile func(r rune) (TileType, error), options ...MapOption) func([]byte) stream.Stream[*Map[TileType]] {
	return func(data []byte) stream.Stream[*Map[TileType]] {
		return &streamParser[TileType]{
			input:     stream.LinesFrom(data),
			options:   options,
			parseTile: parseTile,
		}
	}
}

// Next returns the next map from the input stream.
func (s *streamParser[TileType]) Next() (*Map[TileType], error) {
	for {
		// Read the next line
		line, err := s.input.Next()
		if err != nil {
			if errors.Is(err, io.EOF) && len(s.lines) > 0 {
				m, err := From2DSlices[TileType](s.lines, s.options...)
				if err != nil {
					return nil, err
				}

				s.lines = nil
				return m, nil
			}
			return nil, err
		}

		// If the next line is empty, then loop again
		if line == "" {
			if len(s.lines) > 0 {
				m, err := From2DSlices[TileType](s.lines, s.options...)
				if err != nil {
					return nil, err
				}

				s.lines = nil
				return m, nil
			}
			continue
		}

		if s.lines == nil {
			s.lines = make([][]TileType, 0)
		} else if len(line) != len(s.lines[0]) {
			return nil, errors.Newf("inconsistent line length, expected %d, got %d", len(s.lines[0]), len(line))
		}

		tiles := make([]TileType, len(line))
		for i, r := range line {
			tile, err := s.parseTile(r)
			if err != nil {
				return nil, errors.Wrapf(err, "while parsing tile %d on line %d", i, len(s.lines))
			}

			tiles[i] = tile
		}

		s.lines = append(s.lines, tiles)
	}
}
