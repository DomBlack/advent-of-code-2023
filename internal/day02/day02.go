package day02

import (
	"strconv"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day02 = runner.NewDay(2, parseGames, part1, part2)

type Game struct {
	ID       int
	MaxRed   int
	MaxGreen int
	MaxBlue  int
}

func part1(_ zerolog.Logger, games stream.Stream[Game]) (answer string, err error) {
	// Filter out the games which would have been impossible
	games = stream.Filter(games, func(game Game) (bool, error) {
		return game.MaxRed <= 12 && game.MaxGreen <= 13 && game.MaxBlue <= 14, nil
	})

	gameIDs := stream.Map(games, func(game Game) (int, error) {
		return game.ID, nil
	})

	return stream.SumToString(gameIDs)
}

func part2(_ zerolog.Logger, games stream.Stream[Game]) (answer string, err error) {
	cubePowers := stream.Map(games, func(game Game) (int, error) {
		return game.MaxRed * game.MaxGreen * game.MaxBlue, nil
	})

	return stream.SumToString(cubePowers)
}

func parseGames(input []byte) stream.Stream[Game] {
	lines := stream.LinesFrom(input)

	// Parse each game
	return stream.Map(lines, func(line string) (game Game, err error) {
		// Parse the game ID
		gamePrefix, gamesStr, found := strings.Cut(line, ": ")
		if !found {
			return Game{}, errors.Newf("invalid line: %s", line)
		}

		_, gameIDStr, found := strings.Cut(gamePrefix, " ")
		if !found {
			return Game{}, errors.Newf("invalid line: %s", line)
		}

		game.ID, err = strconv.Atoi(gameIDStr)
		if err != nil {
			return Game{}, errors.Wrapf(err, "invalid line: %s", line)
		}

		rounds := strings.Split(gamesStr, "; ")
		if len(rounds) == 0 {
			return Game{}, errors.Newf("invalid line - no rounds: %s", line)
		}

		// Parse the rounds
		for _, roundStr := range rounds {
			// Parse the round
			round := strings.Split(roundStr, ", ")
			if len(round) == 0 {
				return Game{}, errors.Newf("invalid line - no rounds: %s", line)
			}

			// Parse the colours
			for _, colourStr := range round {
				colour := strings.Split(colourStr, " ")
				if len(colour) != 2 {
					return Game{}, errors.Newf("invalid line - invalid colour: %s", line)
				}

				// Parse the colour
				amount, err := strconv.Atoi(colour[0])
				if err != nil {
					return Game{}, errors.Wrapf(err, "invalid line - invalid colour: %s", line)
				}

				switch colour[1] {
				case "red":
					if amount > game.MaxRed {
						game.MaxRed = amount
					}
				case "green":
					if amount > game.MaxGreen {
						game.MaxGreen = amount
					}
				case "blue":
					if amount > game.MaxBlue {
						game.MaxBlue = amount
					}

				default:
					return Game{}, errors.Newf("invalid line - invalid colour: %s", line)
				}
			}
		}
		return game, nil
	})
}
