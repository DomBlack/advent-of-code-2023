package day09

import (
	"strconv"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day09 = runner.NewStreamingDay(9, parseHistory, part1, part2).
	WithExpectedAnswers("2098530125", "1016")

func part1(log zerolog.Logger, input stream.Stream[History]) (answer string, err error) {
	return stream.SumToString(stream.Map(input, History.NextValue))
}

func part2(log zerolog.Logger, input stream.Stream[History]) (answer string, err error) {
	return stream.SumToString(stream.Map(input, History.PreviousValue))
}

type History []int

func (h History) NextValue() (int, error) {
	// Build the rows
	var rows [][]int
	rows = append(rows, h)
	for {
		previousRow := rows[len(rows)-1]

		allZero := true
		nextRow := make([]int, len(previousRow)-1, len(previousRow))
		for i := 1; i < len(previousRow); i++ {
			nextRow[i-1] = previousRow[i] - previousRow[i-1]
			if nextRow[i-1] != 0 {
				allZero = false
			}
		}
		rows = append(rows, nextRow)
		if allZero {
			break
		}
	}

	for i := len(rows) - 1; i > 0; i-- {
		diff := rows[i][len(rows[i])-1]
		rows[i-1] = append(rows[i-1], rows[i-1][len(rows[i-1])-1]+diff)
	}

	return rows[0][len(rows[0])-1], nil
}

func (h History) PreviousValue() (int, error) {
	// Build the rows
	var rows [][]int
	rows = append(rows, h)
	for {
		previousRow := rows[len(rows)-1]

		allZero := true
		nextRow := make([]int, len(previousRow)-1, len(previousRow))
		for i := 1; i < len(previousRow); i++ {
			nextRow[i-1] = previousRow[i] - previousRow[i-1]
			if nextRow[i-1] != 0 {
				allZero = false
			}
		}
		rows = append(rows, nextRow)
		if allZero {
			break
		}
	}

	for i := len(rows) - 1; i > 0; i-- {
		diff := rows[i][len(rows[i])-1]
		rows[i-1] = append(rows[i-1], rows[i-1][0]-diff)
	}

	return rows[0][len(rows[0])-1], nil
}

func parseHistory(input []byte) stream.Stream[History] {
	return stream.Map(stream.LinesFrom(input), func(line string) (History, error) {
		numStrs := strings.Fields(line)

		nums := make(History, len(numStrs))
		for i, numStr := range numStrs {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse %q as int", numStr)
			}

			nums[i] = num
		}

		return nums, nil
	})
}
