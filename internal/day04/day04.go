package day04

import (
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day04 = runner.NewDay(4, parseCards, part1, part2)

type Card struct {
	Num            int
	Copies         int
	WinningNumbers []int
	PlayedNumbers  []int
}

func parseCards(input []byte) stream.Stream[*Card] {
	return stream.Map(stream.LinesFrom(input), func(line string) (card *Card, err error) {
		if !strings.HasPrefix(line, "Card ") {
			return nil, errors.Newf("line missing prefix: %q", line)
		}
		line = line[5:]

		cardNumStr, rest, found := strings.Cut(line, ":")
		if !found {
			return nil, errors.Newf("unable to find colon: %q", line)
		}

		card = &Card{Copies: 1}
		card.Num, err = strconv.Atoi(strings.TrimSpace(cardNumStr))
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse card number from %q", cardNumStr)
		}

		winningNumbers, playedNumbers, found := strings.Cut(rest, " | ")
		if !found {
			return nil, errors.Wrapf(err, "unable to find pipe separator from %q", rest)
		}

		card.WinningNumbers, err = numbersToSortedSlice(winningNumbers)
		if err != nil {
			return nil, err
		}

		card.PlayedNumbers, err = numbersToSortedSlice(playedNumbers)
		if err != nil {
			return nil, err
		}

		return card, nil
	})
}

func numbersToSortedSlice(str string) ([]int, error) {
	numStrs := strings.Fields(str)
	rtn := make([]int, 0, len(numStrs))

	for _, numStr := range numStrs {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse number from %q", num)
		}
		rtn = append(rtn, num)
	}

	// Sort the slice
	slices.Sort(rtn)

	return rtn, nil
}

func part1(log zerolog.Logger, input stream.Stream[*Card]) (answer string, err error) {
	cardPoints := stream.Map(input, func(card *Card) (int, error) {
		points := 0
		playedIdx := 0

		for _, winning := range card.WinningNumbers {
			// If the winning number is bigger than the number we played
			// then check the next played number
			for winning > card.PlayedNumbers[playedIdx] {
				playedIdx++

				// If we tried all the played numbers, we can quit
				if len(card.PlayedNumbers) <= playedIdx {
					return points, nil
				}
			}

			switch {
			case winning == card.PlayedNumbers[playedIdx]:
				if points == 0 {
					points = 1
				} else {
					points *= 2
				}
			}
		}

		return points, nil
	})

	return stream.SumToString(cardPoints)
}

func part2(log zerolog.Logger, input stream.Stream[*Card]) (answer string, err error) {
	resettable := stream.Resettable(input)

	cardPoints := stream.Map(resettable, func(card *Card) (int, error) {
		winningNumbers := 0
		playedIdx := 0

	winCheck:
		for _, winning := range card.WinningNumbers {
			// If the winning number is bigger than the number we played
			// then check the next played number
			for winning > card.PlayedNumbers[playedIdx] {
				playedIdx++

				// If we tried all the played numbers, we can quit
				if len(card.PlayedNumbers) <= playedIdx {
					break winCheck
				}
			}

			switch {
			case winning == card.PlayedNumbers[playedIdx]:
				winningNumbers++
			}
		}

		log.Debug().Int("card", card.Num).Int("copies", card.Copies).Int("winning_nums", winningNumbers).Msg("Card")
		if winningNumbers > 0 {
			// Update the future cards with additional copies
			resettable.Save()
			for i := 0; i < winningNumbers; i++ {
				toCopy, err := resettable.Next()
				if err != nil {
					if !errors.Is(err, io.EOF) {
						return 0, errors.Wrap(err, "unable to read ahead")
					}
				}

				// For every 1 of this card, the future card gets a new copy too
				toCopy.Copies += card.Copies
			}
			resettable.Restore()
		}

		return card.Copies, nil
	})

	return stream.SumToString(cardPoints)
}
