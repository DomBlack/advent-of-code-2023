package day13

import (
	"io"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day13 = runner.NewStreamingDay(13, parsePatterns, part1, part2).
	WithExpectedAnswers(37975, 32497)

func part1(log zerolog.Logger, input stream.Stream[Pattern]) (answer int, err error) {
	answers := stream.Map(input, func(p Pattern) (int, error) { return p.SummarizeReflection(0) })
	return stream.Sum(answers)
}

func part2(log zerolog.Logger, input stream.Stream[Pattern]) (answer int, err error) {
	answers := stream.Map(input, func(p Pattern) (int, error) { return p.SummarizeReflection(1) })
	return stream.Sum(answers)
}

type Pattern struct {
	Rows []int32 // represents a row of the pattern
	Cols []int32 // represents a column of the pattern
}

// SummarizeReflection returns the number of columns or rows before the
// reflection in the pattern. If there is no reflection, then -1 is returned.
//
// If the reflection is a horizontal reflection, then the number of rows
// before the reflection is multiplied by 100
func (p Pattern) SummarizeReflection(expectedBitFlips int) (int, error) {
	columnsBefore, found := findReflection(p.Cols, expectedBitFlips)
	if found {
		return columnsBefore, nil
	}

	rowsBefore, found := findReflection(p.Rows, expectedBitFlips)
	if found {
		return rowsBefore * 100, nil
	}

	return 0, errors.Newf("No reflection found in pattern:\n%s\n%#v", p, p)
}

func findReflection(pane []int32, expectedBitFlips int) (numBefore int, found bool) {
next:
	for idx := 0; idx < len(pane)-1; idx++ {
		// Check if we've found a possible reflection
		if same, flipsLeft := isSame(pane[idx], pane[idx+1], expectedBitFlips); same {
			for i := 1; idx-i >= 0 && idx+1+i < len(pane); i++ {
				if same, flipsLeft = isSame(pane[idx-i], pane[idx+1+i], flipsLeft); !same {
					// If the reflection doesn't match, then we don't have a reflection
					continue next
				}
			}

			if flipsLeft > 0 {
				// If we have flips left, then we don't have a reflection
				continue next
			}

			// If we get here, then we have a match
			return idx + 1, true
		}
	}

	return 0, false
}

func isSame(a, b int32, allowedBitFlips int) (same bool, flipsLeft int) {
	if a == b {
		return true, allowedBitFlips
	} else if allowedBitFlips == 0 {
		return false, allowedBitFlips
	}

	// If we have a bit flip, then we need to check if it's allowed
	// We can only allow one bit flip, so we need to check if the bit flip
	// is in the same position in both numbers
	for i := 0; i < 32; i++ {
		if a&(1<<i) != b&(1<<i) {
			// If the bit flip is in the same position, then we can allow it
			return isSame(a^(1<<i), b, allowedBitFlips-1)
		}
	}

	return false, allowedBitFlips
}

func (p Pattern) String() string {
	var sb strings.Builder

	for _, row := range p.Rows {
		for i := 0; i < len(p.Cols); i++ {
			if row&(1<<i) > 0 {
				// If the bit is set, write a # as we have a rock
				sb.WriteRune('#')
			} else {
				// Otherwise write a . as we have dust
				sb.WriteRune('.')
			}
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func parsePatterns(input []byte) stream.Stream[Pattern] {
	return &patternParser{
		input: stream.LinesFrom(input),
	}
}

type patternParser struct {
	input   stream.Stream[string]
	current Pattern
}

func (p *patternParser) Next() (Pattern, error) {
	for {
		// Read the next line
		line, err := p.input.Next()
		if err != nil {
			// If we've reached EOF on the source stream and we have a pattern
			// return it, otherwise return the error
			if errors.Is(err, io.EOF) && len(p.current.Rows) > 0 {
				toReturn := p.current
				p.current = Pattern{}
				return toReturn, nil
			} else {
				return Pattern{}, err
			}
		}

		// If we hit an empty line, return the current pattern
		if line == "" {
			toReturn := p.current
			p.current = Pattern{}
			return toReturn, nil
		}

		// Sanity check the width of the line
		if len(p.current.Cols) > 0 && len(p.current.Cols) != len(line) {
			return Pattern{}, errors.Newf("Line %s is not the same width as previous lines in this pattern, expected width %d", line, len(p.current.Cols))
		} else if len(p.current.Cols) == 0 {
			p.current.Cols = make([]int32, len(line))
		}

		row := int32(0)

		for col, c := range line {
			switch c {
			case '.':
				// no-op
			case '#':
				// Record the rock on this row
				row |= 1 << col

				// And then on the column
				p.current.Cols[col] |= 1 << len(p.current.Rows)

			default:
				return Pattern{}, errors.Newf("Unknown character %c", c)
			}
		}

		p.current.Rows = append(p.current.Rows, row)
	}
}
