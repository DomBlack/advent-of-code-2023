package day14

import (
	"fmt"
	"image/color"

	"github.com/DomBlack/advent-of-code-2023/pkg/maps"
	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var (
	Day14 = runner.NewDay(14, parseFunc, part1, part2).
		WithExpectedAnswers(105982, 85175)

	parseFunc = maps.NewParseFunc(func(r rune) (Rocks, error) {
		switch r {
		case '.':
			return Empty, nil
		case 'O':
			return Rounded, nil
		case '#':
			return Cubed, nil
		default:
			return Empty, errors.Newf("invalid rock type: %v", r)
		}
	})
)

func part1(log zerolog.Logger, input *maps.Map[Rocks]) (answer int, err error) {
	maps.Tilt(input, maps.North, Rounded)

	return load(input), nil
}

func part2(log zerolog.Logger, input *maps.Map[Rocks]) (answer int, err error) {
	const spinCount = 1_000_000_000
	cache := make(map[string]int)

	spinCycle := func(i int) {
		// Run the spin cycle (North, West, South, East)
		maps.Tilt(input, maps.North, Rounded)
		maps.Tilt(input, maps.West, Rounded)
		maps.Tilt(input, maps.South, Rounded)
		maps.Tilt(input, maps.East, Rounded)
		input.CaptureFrame(fmt.Sprintf("Spin Cycle: %d", i), 1)
	}

	file := runner.Output(14)
	if file != "" {
		input.StartCapturingFrames()
	}

	// Run the spin cycle 1 billion times or until we find a loop
	var cacheKey string
	i := 1
	for i < spinCount {
		spinCycle(i)

		cacheKey = input.String()
		if _, ok := cache[cacheKey]; ok {
			break
		}
		cache[cacheKey] = i
		i++
	}

	loopStart := cache[cacheKey]
	loopLength := i - loopStart
	if loopLength == 0 {
		return 0, errors.New("no loop found")
	}
	log.Debug().Int("loop_start", loopStart).Int("loop_length", loopLength).Int("current_idx", i).Msg("found loop")

	// Fast forward to the end
	startAt := spinCount - ((spinCount - loopStart) % loopLength)
	log.Debug().Int("iteration", startAt).Msg("fast forwarding to end")
	for i = startAt; i < spinCount; i++ {
		spinCycle(i)
	}

	answer = load(input)

	// Save the animation
	if file != "" {
		input.StopCapturingFrames(fmt.Sprintf("Spin Cycle: %d - Answer: %d", i, answer))
		if err := input.SaveAnimationGIF(file + ".gif"); err != nil {
			return 0, errors.Wrap(err, "failed to save animation gif")
		}
	}

	return answer, nil
}

func load(m *maps.Map[Rocks]) (sum int) {
	for idx, tile := range m.Tiles {
		if tile == Rounded {
			_, y := m.PositionOf(idx)
			sum += m.Height - y
		}
	}

	return sum
}

type Rocks uint8

const (
	Empty Rocks = iota
	Rounded
	Cubed

	eof
)

func (r Rocks) Valid() bool {
	return r < eof
}

func (r Rocks) Rune() rune {
	switch r {
	case Empty:
		return '.'
	case Rounded:
		return 'O'
	case Cubed:
		return '#'
	default:
		panic(fmt.Sprintf("invalid rock type: %v", r))
	}
}

func (r Rocks) Colour() color.Color {
	switch r {
	case Empty:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // White
	case Rounded:
		return color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green
	case Cubed:
		return color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black
	default:
		panic(fmt.Sprintf("invalid rock type: %v", r))
	}
}
