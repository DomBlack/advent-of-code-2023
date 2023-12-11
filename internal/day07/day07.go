package day07

import (
	"slices"
	"strconv"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day07 = runner.NewStreamingDay(7, parseHands,
	solver(func(hand Hand) (Hand, error) { return hand, nil }),
	solver(withJokerCards),
).
	WithExpectedAnswers("250453939", "248652697")

func solver(mapper func(Hand) (Hand, error)) func(log zerolog.Logger, input stream.Stream[Hand]) (answer string, err error) {
	return func(log zerolog.Logger, input stream.Stream[Hand]) (answer string, err error) {
		hands, err := stream.Collect(stream.Map(input, mapper))
		if err != nil {
			return "", err
		}

		slices.SortFunc(hands, func(a, b Hand) int {
			typeDiff := int(a.Type) - int(b.Type)
			if typeDiff != 0 {
				return typeDiff
			}

			for i := 0; i < 5; i++ {
				diff := int(a.Cards[i]) - int(b.Cards[i])
				if diff != 0 {
					return diff
				}
			}

			return 0
		})

		winnings := 0
		for i, hand := range hands {
			log.Debug().Int("i", i).Int("type", int(hand.Type)).Msg(hand.String())
			winnings += hand.Bid * (i + 1)
		}

		return strconv.Itoa(winnings), nil
	}
}

type HandType uint8

const (
	HighCard HandType = iota
	Pair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type Hand struct {
	Cards [5]uint8
	Bid   int
	Type  HandType
}

func (h Hand) String() string {
	var sb strings.Builder

	for _, card := range h.Cards {
		switch card {
		case 10:
			sb.WriteString("T")
		case 11, 1:
			sb.WriteString("J")
		case 12:
			sb.WriteString("Q")
		case 13:
			sb.WriteString("K")
		case 14:
			sb.WriteString("A")
		default:
			sb.WriteByte(card + '0')
		}
	}
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(h.Bid))

	sb.WriteByte(' ')
	switch h.Type {
	case HighCard:
		sb.WriteString("HighCard")
	case Pair:
		sb.WriteString("Pair")
	case TwoPair:
		sb.WriteString("TwoPair")
	case ThreeOfAKind:
		sb.WriteString("ThreeOfAKind")
	case FullHouse:
		sb.WriteString("FullHouse")
	case FourOfAKind:
		sb.WriteString("FourOfAKind")
	case FiveOfAKind:
		sb.WriteString("FiveOfAKind")
	}

	return sb.String()
}

func parseHands(bytes []byte) stream.Stream[Hand] {
	return stream.Map(stream.LinesFrom(bytes), func(line string) (Hand, error) {
		cardsStr, bidStr, found := strings.Cut(line, " ")
		if !found {
			return Hand{}, errors.Newf("Could not find bid in line: %q", line)
		}

		bid, err := strconv.Atoi(bidStr)
		if err != nil {
			return Hand{}, errors.Wrapf(err, "Could not parse bid: %q", bid)
		}

		groupings := map[uint8]uint8{}

		cards := [5]uint8{}
		for i := 0; i < 5; i++ {
			switch cardsStr[i] {
			case '2', '3', '4', '5', '6', '7', '8', '9':
				cards[i] = cardsStr[i] - '0'
			case 'T':
				cards[i] = 10
			case 'J':
				cards[i] = 11
			case 'Q':
				cards[i] = 12
			case 'K':
				cards[i] = 13
			case 'A':
				cards[i] = 14
			}

			groupings[cards[i]]++
		}

		maxCount := uint8(0)
		minCount := uint8(255)
		twoPairCount := 0

		for _, count := range groupings {
			if count > maxCount {
				maxCount = count
			}

			if count < minCount {
				minCount = count
			}

			if count == 2 {
				twoPairCount++
			}
		}

		handType := HighCard
		switch {
		case maxCount == 5:
			handType = FiveOfAKind
		case maxCount == 4:
			handType = FourOfAKind
		case maxCount == 3 && minCount == 2:
			handType = FullHouse
		case maxCount == 3:
			handType = ThreeOfAKind
		case maxCount == 2 && twoPairCount == 2:
			handType = TwoPair
		case maxCount == 2:
			handType = Pair
		}

		return Hand{
			Cards: cards,
			Bid:   bid,
			Type:  handType,
		}, nil
	})
}

func withJokerCards(hand Hand) (Hand, error) {
	numJokers := 0
	groupings := map[uint8]uint8{}

	for i := 0; i < 5; i++ {
		if hand.Cards[i] == 11 {
			hand.Cards[i] = 1
			numJokers++
		} else {
			groupings[hand.Cards[i]]++
		}
	}

	maxCount := uint8(0)
	minCount := uint8(255)
	twoPairCount := 0

	for _, count := range groupings {
		if count > maxCount {
			maxCount = count
		}

		if count < minCount {
			minCount = count
		}

		if count == 2 {
			twoPairCount++
		}
	}

	if numJokers > 0 {
		switch {
		case maxCount == 4:
			switch numJokers {
			case 1:
				hand.Type = FiveOfAKind
			default:
				return hand, errors.Newf("Unexpected numJokers: %d", numJokers)
			}
		case maxCount == 3:
			switch numJokers {
			case 1:
				hand.Type = FourOfAKind
			case 2:
				hand.Type = FiveOfAKind
			default:
				return hand, errors.Newf("Unexpected numJokers: %d", numJokers)
			}
		case maxCount == 2 && twoPairCount == 2:
			switch numJokers {
			case 1:
				hand.Type = FullHouse
			case 2:
				hand.Type = FourOfAKind
			case 3:
				hand.Type = FiveOfAKind
			default:
				return hand, errors.Newf("Unexpected numJokers: %d", numJokers)
			}
		case maxCount == 2:
			switch numJokers {
			case 1:
				hand.Type = ThreeOfAKind
			case 2:
				hand.Type = FourOfAKind
			case 3:
				hand.Type = FiveOfAKind
			default:
				return hand, errors.Newf("Unexpected numJokers: %d", numJokers)
			}
		case maxCount == 1:
			switch numJokers {
			case 1:
				hand.Type = Pair
			case 2:
				hand.Type = ThreeOfAKind
			case 3:
				hand.Type = FourOfAKind
			case 4:
				hand.Type = FiveOfAKind
			default:
				return hand, errors.Newf("Unexpected numJokers: %d", numJokers)
			}

		case maxCount == 0:
			hand.Type = FiveOfAKind

		default:
			return hand, errors.Newf("Unexpected maxCount: %d", maxCount)
		}
	}

	return hand, nil
}
