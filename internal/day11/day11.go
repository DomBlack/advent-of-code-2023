package day11

import (
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/vectors/vec2"
	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day11 = runner.NewDay(11, parseUniverse, part1, part2).
	WithExpectedAnswers(9724940, 569052586852)

func part1(log zerolog.Logger, input Universe) (answer int, err error) {
	expanded := input.expand(2)

	return expanded.sumOfLengths(log), nil
}

func part2(log zerolog.Logger, input Universe) (answer int, err error) {
	expanded := input.expand(1000000)

	return expanded.sumOfLengths(log), nil
}

type Universe []Galaxy

type Galaxy struct {
	Position vec2.Vec2i
}

func parseUniverse(input []byte) (Universe, error) {
	rtn := make(Universe, 0)

	// Parse the initial galaxies
	x := 0
	y := 0

	for _, r := range strings.TrimSpace(string(input)) {

		switch r {
		case '\n':
			y++
			x = 0
			continue
		case '.':
			// no-op
		case '#':
			rtn = append(rtn, Galaxy{vec2.Vec2i{x, y}})
		default:
			return nil, errors.Newf("Unknown character: %c", r)
		}

		x++
	}

	return rtn, nil
}

// sumOfLengths returns the sum of the lengths between all the galaxies
func (u Universe) sumOfLengths(log zerolog.Logger) int {
	sum := 0

	for i, galaxy := range u {
		log.Debug().Int("galaxy", i+1).Int("x", galaxy.Position[0]).Int("y", galaxy.Position[1]).Msg("Galaxy")

		for j := i + 1; j < len(u); j++ {
			galaxy2 := u[j]
			length := galaxy2.Position.Sub(galaxy.Position).Length()
			log.Debug().Int("galaxy", i+1).Int("galaxy2", j+1).Int("length", length).Msg("Distance")

			sum += length
		}
	}

	return sum
}

// expand expands the universe by the given amount
//
// A growth factor of 2 will double the size of any empty rows/columns (i.e. 1 column becomes 2 columns)
// A growth factor of 1000 will increase the size of any empty rows/columns by 1000-1 (i.e. 1 column becomes 1000 columns)
func (u Universe) expand(growthFactor int) Universe {
	growthFactor -= 1

	rtn := make(Universe, len(u))

	// Copy the universe
	copy(rtn, u)

	// Mark all the rows and columns that have a galaxy
	columnsWithGalaxy := make(map[int]struct{})
	rowsWithGalaxy := make(map[int]struct{})
	width := 0
	height := 0
	for _, galaxy := range rtn {
		columnsWithGalaxy[galaxy.Position[0]] = struct{}{}
		rowsWithGalaxy[galaxy.Position[1]] = struct{}{}

		if galaxy.Position[0] > width {
			width = galaxy.Position[0]
		}
		if galaxy.Position[1] > height {
			height = galaxy.Position[1]
		}
	}

	// Starting with the higest column
	for x := width - 1; x >= 0; x-- {
		// check if there's no galaxy in that column
		if _, ok := columnsWithGalaxy[x]; ok {
			continue
		}

		// and if there isn't then move all the galaxies to the right of it by 1 column
		for i := 0; i < len(rtn); i++ {
			if rtn[i].Position[0] >= x {
				rtn[i].Position[0] += growthFactor
			}
		}
	}

	// now repeat for the rows
	for y := height - 1; y >= 0; y-- {
		// check if there's no galaxy in that row
		if _, ok := rowsWithGalaxy[y]; ok {
			continue
		}

		// and if there isn't then move all the galaxies to the below it by 1 row
		for i := 0; i < len(rtn); i++ {
			if rtn[i].Position[1] >= y {
				rtn[i].Position[1] += growthFactor
			}
		}
	}

	return rtn
}
