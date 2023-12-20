package day17

import (
	"fmt"
	"image/color"

	"github.com/DomBlack/advent-of-code-2023/pkg/algorithms/astar"
	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/vectors/vec2"
	"github.com/DomBlack/advent-of-code-2023/pkg/maps"
	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var (
	Day17 = runner.NewDay(17, parseFunc, part1, part2).
		WithExpectedAnswers(907, 1057)

	// colorPalette starts with a faint red and goes to a deep red in 9 steps
	colorPalette = []color.Color{
		color.RGBA{0xff, 0xff, 0xff, 0xff}, // base is white
		color.RGBA{0xff, 0x88, 0x88, 0xff},
		color.RGBA{0xff, 0x77, 0x77, 0xff},
		color.RGBA{0xff, 0x66, 0x66, 0xff},
		color.RGBA{0xff, 0x55, 0x55, 0xff},
		color.RGBA{0xff, 0x44, 0x44, 0xff},
		color.RGBA{0xff, 0x33, 0x33, 0xff},
		color.RGBA{0xff, 0x22, 0x22, 0xff},
		color.RGBA{0xff, 0x11, 0x11, 0xff},
		color.RGBA{0xff, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff}, // green for the path
		color.RGBA{0xdd, 0xff, 0xdd, 0xff}, // light green for path options

	}

	parseFunc = maps.NewParseFunc(func(r rune) (Tile, error) {
		switch r {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return Tile(r - '1'), nil
		default:
			return 0, errors.Newf("invalid tile type: %v", r)
		}
	}, maps.WithColourPalette(colorPalette))
)

type SearchState struct {
	pos  maps.Pos
	Dir  vec2.Vec2
	Dist int
}

func (s SearchState) Pos() maps.Pos {
	return s.pos
}

func part1(ctx *runner.Context, log zerolog.Logger, input *maps.Map[Tile]) (answer int, err error) {
	// Compute the average tile cost
	sumOfAllTiles := 0
	for _, tile := range input.Tiles {
		sumOfAllTiles += int(tile)
	}
	avgTileCost := sumOfAllTiles / len(input.Tiles)

	return findPath(ctx, input, 1, 3, avgTileCost)
}

func part2(ctx *runner.Context, log zerolog.Logger, input *maps.Map[Tile]) (answer int, err error) {
	return findPath(ctx, input, 4, 10, 2)
}

type Tile uint8

const (
	Loss1 Tile = iota
	Loss2
	Loss3
	Loss4
	Loss5
	Loss6
	Loss7
	Loss8
	Loss9
	loss_bits = 15

	Path       = 16
	PathOption = 32
)

func (t Tile) Valid() bool {
	return t <= Loss9
}

func (t Tile) Rune() rune {
	return rune((t & 15) + '0')
}

func (t Tile) Cost() int {
	return int(t&15) + 1
}

func (t Tile) Colour() color.Color {
	switch {
	case t&Path == Path:
		return colorPalette[10]
	case t&PathOption == PathOption:
		return colorPalette[11]
	default:
		return colorPalette[t&loss_bits+1]
	}
}

func findPath(ctx *runner.Context, input *maps.Map[Tile], minDist, maxDist int, avgCostPerTile int) (cost int, err error) {
	// Clean out our state
	for i := range input.Tiles {
		input.Tiles[i] = input.Tiles[i] &^ Path &^ PathOption
	}

	input.StartCapturingFrames(ctx)

	// Then run A* to find the path
	goal := maps.Pos{input.Width - 1, input.Height - 1}
	cost, path, err := astar.Search(
		input,
		SearchState{pos: maps.Pos{0, 0}},
		func(from SearchState) bool { return from.pos == goal && from.Dist >= minDist },
		neighbourFunc(input, minDist, maxDist),
		func(from SearchState) int { return goal.Sub(from.pos).Length() * avgCostPerTile },
		Path, PathOption,
	)

	if err != nil {
		return 0, errors.Wrap(err, "failed to find path")
	}

	for _, pos := range path {
		input.AddFlagAt(pos, Path)
	}
	input.StopCapturingFrames(fmt.Sprintf("Path Length: %d - Cost: %d", len(path), cost))

	if err := input.SaveAnimationGIF(ctx); err != nil {
		return 0, errors.Wrap(err, "failed to save animation")
	}

	return cost, nil
}

func neighbourFunc(input *maps.Map[Tile], minDist, maxDist int) func(from SearchState) (rtn []SearchState) {
	return func(from SearchState) (rtn []SearchState) {
		// We must move in the same direction for at least minDist blocks
		if from.Dist < minDist && from.Dir != vec2.Zero {
			newPos := from.pos.Add(from.Dir)
			if !input.InBounds(newPos) {
				return nil
			}
			return []SearchState{
				{
					pos:  newPos,
					Dir:  from.Dir,
					Dist: from.Dist + 1,
				},
			}
		}

		for _, neighbour := range input.Neighbours(from.pos) {
			newDir := neighbour.Sub(from.pos)

			if newDir == from.Dir.Neg() {
				// Prevent going backwards
				continue
			}

			newDist := 1
			if newDir == from.Dir {
				newDist = from.Dist + 1

				// Prevent going over the maxDist in the same direction
				if newDist > maxDist {
					continue
				}
			}

			rtn = append(rtn, SearchState{
				pos:  neighbour,
				Dir:  newDir,
				Dist: newDist,
			})
		}
		return rtn
	}
}
