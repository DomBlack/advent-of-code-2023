package day01

import (
	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/trie"
	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day01 = runner.NewStreamingDay(1, stream.LinesFrom, implementation(part1Tree), implementation(part2Tree)).
	WithExpectedAnswers(54630, 54770)

var part1Tree = trie.New[int]().
	MustInsert("0", 0).
	MustInsert("1", 1).
	MustInsert("2", 2).
	MustInsert("3", 3).
	MustInsert("4", 4).
	MustInsert("5", 5).
	MustInsert("6", 6).
	MustInsert("7", 7).
	MustInsert("8", 8).
	MustInsert("9", 9)

var part2Tree = trie.New[int]().
	MustInsert("0", 0).
	MustInsert("1", 1).
	MustInsert("2", 2).
	MustInsert("3", 3).
	MustInsert("4", 4).
	MustInsert("5", 5).
	MustInsert("6", 6).
	MustInsert("7", 7).
	MustInsert("8", 8).
	MustInsert("9", 9).
	MustInsert("one", 1).
	MustInsert("two", 2).
	MustInsert("three", 3).
	MustInsert("four", 4).
	MustInsert("five", 5).
	MustInsert("six", 6).
	MustInsert("seven", 7).
	MustInsert("eight", 8).
	MustInsert("nine", 9)

func implementation(tree *trie.Trie[int]) func(log zerolog.Logger, input stream.Stream[string]) (answer int, err error) {
	return func(log zerolog.Logger, input stream.Stream[string]) (answer int, err error) {
		digitsPerLine := stream.Map(input, func(line string) (int, error) {
			matches := tree.SubstrMatches(line)

			if len(matches) <= 0 {
				return 0, errors.New("no digits found on line")
			}
			first := matches[0]
			last := matches[len(matches)-1]

			log.Debug().Int("first", first).Int("last", last).Str("line", line).Msg("Found digits")

			return first*10 + last, nil
		})

		return stream.Sum(digitsPerLine)
	}
}
