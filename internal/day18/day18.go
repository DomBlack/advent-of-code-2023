package day18

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/algorithms/floodfill"
	"github.com/DomBlack/advent-of-code-2023/pkg/algorithms/polygonarea"
	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/vectors/vec2"
	"github.com/DomBlack/advent-of-code-2023/pkg/maps"
	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day18 = runner.NewDay(18, parseInstructions, part1, part2)

// Instruction is a single instruction for the dig plan
type Instruction struct {
	Direction maps.Pos
	Length    int

	SwappedDirection maps.Pos
	SwappedLength    int
}

func (i Instruction) Swap() Instruction {
	return Instruction{
		Direction:        i.SwappedDirection,
		Length:           i.SwappedLength,
		SwappedDirection: i.Direction,
		SwappedLength:    i.Length,
	}
}

func part1(ctx *runner.Context, log zerolog.Logger, input []Instruction) (answer int, err error) {
	m, startPos := convertToMap(ctx, input)

	floodfill.Fill(m, startPos.Add(maps.Pos{1, 1}), Lava)

	count := 0
	for _, tile := range m.Tiles {
		if tile != Empty {
			count++
		}
	}

	m.StopCapturingFrames(fmt.Sprintf("Answer: %d", count))
	if err := m.SaveAnimationGIF(ctx); err != nil {
		return 0, errors.Wrap(err, "failed to write animation")
	}
	return count, nil
}

func part2(ctx *runner.Context, log zerolog.Logger, input []Instruction) (answer int, err error) {
	vertexes := make([]vec2.Vec2, 0, len(input))
	var pos maps.Pos
	vertexes = append(vertexes, pos)

	for _, instruction := range input {
		for i := 0; i < instruction.SwappedLength; i++ {
			pos = pos.Add(instruction.SwappedDirection)
			vertexes = append(vertexes, pos)
		}
	}
	polygonArea := polygonarea.Area(vertexes)

	return polygonArea, nil
}

func parseInstructions(input []byte) ([]Instruction, error) {
	return stream.Collect(stream.Map(stream.LinesFrom(input), func(line string) (rtn Instruction, err error) {
		switch line[0] {
		// Parse the direction
		case 'R':
			rtn.Direction = maps.Pos{1, 0}
		case 'L':
			rtn.Direction = maps.Pos{-1, 0}
		case 'U':
			rtn.Direction = maps.Pos{0, -1}
		case 'D':
			rtn.Direction = maps.Pos{0, 1}
		default:
			return rtn, errors.Newf("invalid direction: %c", line[0])
		}
		if line[1] != ' ' {
			return rtn, errors.Newf("invalid line: %q", line)
		}

		// Parse the length and update the direction
		lengthStr, hexCode, found := strings.Cut(line[2:], " ")
		if !found {
			return rtn, errors.Newf("invalid line: %q", line)
		}
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return rtn, errors.Wrap(err, "invalid length")
		}
		rtn.Length = length

		// Parse the "colour" (really the swapped instructions)
		if len(hexCode) != 9 {
			return rtn, errors.Newf("invalid hex code: %q", hexCode)
		}
		hexCode = hexCode[2:8]

		// The last digit encodes the swapped direction
		switch hexCode[5] {
		case '0':
			rtn.SwappedDirection = maps.Pos{1, 0}
		case '1':
			rtn.SwappedDirection = maps.Pos{0, -1}
		case '2':
			rtn.SwappedDirection = maps.Pos{-1, 0}
		case '3':
			rtn.SwappedDirection = maps.Pos{0, 1}
		default:
			return rtn, errors.Newf("invalid swapped direction: %c", hexCode[7])
		}

		// The remaining 6 digits encode the swapped length
		swappedLength, err := strconv.ParseUint(hexCode[:5], 16, 32)
		if err != nil {
			return rtn, errors.Wrap(err, "invalid swapped length")
		}
		rtn.SwappedLength = int(swappedLength)

		return rtn, nil
	}))
}

func convertToMap(ctx *runner.Context, instructions []Instruction) (*maps.Map[Tile], maps.Pos) {
	// Work out the size of the map
	minX := math.MaxInt
	maxX := math.MinInt
	minY := math.MaxInt
	maxY := math.MinInt

	var pos maps.Pos
	for _, instruction := range instructions {
		pos = pos.Add(instruction.Direction.Scale(instruction.Length))
		minX = min(minX, pos[0])
		maxX = max(maxX, pos[0])
		minY = min(minY, pos[1])
		maxY = max(maxY, pos[1])
	}

	// Create the map
	m := maps.New[Tile](maxX-minX+1, maxY-minY+1)

	m.StartCapturingFrames(ctx)

	pos = maps.Pos{-minX, -minY}
	startPos := pos
	for _, instruction := range instructions {
		for i := 0; i < instruction.Length; i++ {
			m.Set(pos, Trench)
			pos = pos.Add(instruction.Direction)
		}

		m.CaptureFrame("Digging", 1)
	}

	return m, startPos
}

type Tile uint8

const (
	Empty Tile = iota
	Trench
	Lava
)

func (t Tile) Valid() bool {
	return t <= Lava
}

func (t Tile) Rune() rune {
	switch t {
	case Empty:
		return '.'
	case Trench:
		return '#'
	case Lava:
		return '*'
	default:
		panic("invalid tile")
	}
}

func (t Tile) Colour() color.Color {
	switch t {
	case Empty:
		return color.White
	case Trench:
		return color.Black
	case Lava:
		return color.RGBA{R: 200, G: 0, B: 0, A: 255}
	default:
		panic("invalid tile")
	}
}
