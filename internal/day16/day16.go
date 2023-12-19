package day16

import (
	"fmt"
	"image/color"
	"image/draw"
	"strconv"

	"github.com/DomBlack/advent-of-code-2023/pkg/maps"
	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var (
	Day16 = runner.NewDay(16, parseContraption, part1, part2)

	colorPalette = []color.Color{
		color.White,                // Empty
		color.Black,                // Mirrors
		color.RGBA{G: 255, A: 255}, // Laser Beam Head
		color.RGBA{G: 200, A: 255}, // Laser Beam Tail
	}

	parseContraption = maps.NewParseFunc(parseMirrorTile, maps.WithTileRender(mirrorRenderer), maps.WithColourPalette(colorPalette))
)

type laserHead struct {
	x   int
	y   int
	dir Tile
}

func part1(ctx *runner.Context, log zerolog.Logger, input *maps.Map[Tile]) (answer int, err error) {
	// First reset the map
	for i := range input.Tiles {
		input.Tiles[i] &= NonLaserBits
	}

	input.StartCapturingFrames(ctx)

	answer, err = runLasers(input, 0, 0, LaserRight, true)
	if err != nil {
		return 0, errors.Wrap(err, "failed to run lasers")
	}

	// Save the animation
	input.StopCapturingFrames(fmt.Sprintf("Answer: %d energized tiles", answer))
	if err := input.SaveAnimationGIF(ctx); err != nil {
		return 0, errors.Wrap(err, "failed to save animation gif")
	}

	return answer, nil
}

func part2(ctx *runner.Context, log zerolog.Logger, input *maps.Map[Tile]) (answer int, err error) {
	// First reset the map
	for i := range input.Tiles {
		input.Tiles[i] &= NonLaserBits
	}

	input.StartCapturingFrames(ctx)

	bestX := 0
	bestY := 0
	bestDir := LaserRight

	for x := 0; x < input.Width; x++ {
		thisOption, err := runLasers(input, x, 0, LaserDown, false)
		if err != nil {
			return 0, errors.Wrap(err, "failed to run lasers")
		}
		input.CaptureFrame(fmt.Sprintf("(%d, 0) = %d energized tiles", x, thisOption), 1)

		if thisOption > answer {
			bestX = x
			bestY = 0
			bestDir = LaserDown
			answer = thisOption
		}

		thisOption, err = runLasers(input, x, input.Height-1, LaserUp, false)
		if err != nil {
			return 0, errors.Wrap(err, "failed to run lasers")
		}
		input.CaptureFrame(fmt.Sprintf("(%d, %d) = %d energized tiles", x, input.Height-1, thisOption), 1)

		if thisOption > answer {
			bestX = x
			bestY = input.Height - 1
			bestDir = LaserUp
			answer = thisOption
		}
	}

	for y := 0; y < input.Height; y++ {
		thisOption, err := runLasers(input, 0, y, LaserRight, false)
		if err != nil {
			return 0, errors.Wrap(err, "failed to run lasers")
		}
		input.CaptureFrame(fmt.Sprintf("(%d, 0) = %d energized tiles", 0, y), 1)

		if thisOption > answer {
			bestX = 0
			bestY = y
			bestDir = LaserRight
			answer = thisOption
		}

		thisOption, err = runLasers(input, input.Width-1, y, LaserLeft, false)
		if err != nil {
			return 0, errors.Wrap(err, "failed to run lasers")
		}
		input.CaptureFrame(fmt.Sprintf("(%d, %d) = %d energized tiles", input.Width-1, y, thisOption), 1)

		if thisOption > answer {
			bestX = input.Width - 1
			bestY = y
			bestDir = LaserLeft
			answer = thisOption
		}
	}

	// Save the animation with the best layout
	_, _ = runLasers(input, bestX, bestY, bestDir, false)
	input.StopCapturingFrames(fmt.Sprintf("Answer: (%d, %d) = %d energized tiles", bestX, bestY, answer))

	if err := input.SaveAnimationGIF(ctx); err != nil {
		return 0, errors.Wrap(err, "failed to save animation gif")
	}

	return answer, nil
}

func runLasers(input *maps.Map[Tile], startX int, startY int, initDirection Tile, recordEachStep bool) (answer int, err error) {
	// First reset the map
	for i := range input.Tiles {
		input.Tiles[i] &= NonLaserBits
	}

	// Set up the initial state
	laserHeads := make([]laserHead, 0)
	nextLaserHeads := make([]laserHead, 0)
	laserHeads = append(laserHeads, laserHead{startX, startY, initDirection})
	input.AddFlagAt(startX, startY, LaserHead|initDirection)

	step := 1

	var newDirs []Tile

	for len(laserHeads) > 0 {
		nextLaserHeads = nextLaserHeads[:0]

		for _, head := range laserHeads {
			tile, valid := input.Get(head.x, head.y)
			if !valid {
				continue
			}
			input.RemoveFlagAt(head.x, head.y, LaserHead)

			// Calculate the new directions
			switch head.dir {
			case LaserUp:
				if tile&NonLaserBits == HorizontalSplitter {
					newDirs = append(newDirs, LaserLeft, LaserRight)
				} else {
					newDir := LaserUp
					if tile&NonLaserBits == BackslashMirror {
						newDir = LaserLeft
					} else if tile&NonLaserBits == ForwardSlashMirror {
						newDir = LaserRight
					}

					newDirs = append(newDirs, newDir)
				}

			case LaserDown:
				if tile&NonLaserBits == HorizontalSplitter {
					newDirs = append(newDirs, LaserLeft, LaserRight)
				} else {
					newDir := LaserDown
					if tile&NonLaserBits == BackslashMirror {
						newDir = LaserRight
					} else if tile&NonLaserBits == ForwardSlashMirror {
						newDir = LaserLeft
					}

					newDirs = append(newDirs, newDir)
				}

			case LaserLeft:
				if tile&NonLaserBits == VerticalSplitter {
					newDirs = append(newDirs, LaserUp, LaserDown)
				} else {
					newDir := LaserLeft
					if tile&NonLaserBits == BackslashMirror {
						newDir = LaserUp
					} else if tile&NonLaserBits == ForwardSlashMirror {
						newDir = LaserDown
					}

					newDirs = append(newDirs, newDir)
				}

			case LaserRight:
				if tile&NonLaserBits == VerticalSplitter {
					newDirs = append(newDirs, LaserUp, LaserDown)
				} else {
					newDir := LaserRight
					if tile&NonLaserBits == BackslashMirror {
						newDir = LaserDown
					} else if tile&NonLaserBits == ForwardSlashMirror {
						newDir = LaserUp
					}

					newDirs = append(newDirs, newDir)
				}
			default:
				return 0, errors.Newf("unknown laser direction: %d", head.dir)
			}

			// If the tile is empty, then we need to draw the laser beam
			for _, dir := range newDirs {
				newX := head.x
				newY := head.y

				switch dir {
				case LaserUp:
					newY--
				case LaserDown:
					newY++
				case LaserLeft:
					newX--
				case LaserRight:
					newX++
				default:
					return 0, errors.Newf("unknown laser direction: %d", dir)
				}

				// If the new pos isn't valid or is already a laser beam in the same direction, then skip it
				newTile, valid := input.Get(newX, newY)
				if !valid || newTile&dir != 0 {
					continue
				}

				input.AddFlagAt(newX, newY, LaserHead|dir)
				nextLaserHeads = append(nextLaserHeads, laserHead{newX, newY, dir})
			}
			newDirs = newDirs[:0]
		}

		laserHeads, nextLaserHeads = nextLaserHeads, laserHeads
		if recordEachStep {
			input.CaptureFrame("Step "+strconv.Itoa(step), 1)
		}
		step++

		if step > 1_000_000_000 {
			return 0, errors.New("too many steps")
		}
	}

	// Count the number of energized tiles
	for _, tile := range input.Tiles {
		if tile & ^NonLaserBits != 0 {
			answer++
		}
	}

	return answer, nil
}

type Tile uint8

const (
	Empty Tile = iota
	BackslashMirror
	ForwardSlashMirror
	VerticalSplitter
	HorizontalSplitter

	eof

	// Laser tracking flags
	NonLaserBits Tile = 7
	LaserDown    Tile = 2 << 2
	LaserUp      Tile = 2 << 3
	LaserLeft    Tile = 2 << 4
	LaserRight   Tile = 2 << 5
	LaserHead    Tile = 2 << 6
)

func parseMirrorTile(r rune) (Tile, error) {
	switch r {
	case '.':
		return Empty, nil
	case '\\':
		return BackslashMirror, nil
	case '/':
		return ForwardSlashMirror, nil
	case '|':
		return VerticalSplitter, nil
	case '-':
		return HorizontalSplitter, nil
	default:
		return Empty, errors.Newf("unknown tile type: %c", r)
	}
}

func (t Tile) Valid() bool {
	return t < eof
}

func (t Tile) Rune() rune {
	switch t & NonLaserBits {
	case Empty:
		return '.'
	case BackslashMirror:
		return '\\'
	case ForwardSlashMirror:
		return '/'
	case VerticalSplitter:
		return '|'
	case HorizontalSplitter:
		return '-'
	default:
		panic("unknown tile type")
	}
}

func (t Tile) Colour() color.Color {
	if t & ^NonLaserBits != 0 {
		// Energized tiles are green
		return colorPalette[2]
	}

	switch t & NonLaserBits {
	case Empty:
		return color.White
	case BackslashMirror:
		return color.Black
	case ForwardSlashMirror:
		return color.Black
	case VerticalSplitter:
		return color.Black
	case HorizontalSplitter:
		return color.Black
	default:
		panic("unknown tile type")
	}
}

func mirrorRenderer(tile maps.AnyTile, img draw.Image, x int, y int, size int) {
	t := tile.(Tile)

	laserColor := colorPalette[3]
	if t&LaserHead != 0 {
		laserColor = colorPalette[2]
	}

	switch t & NonLaserBits {
	case Empty:
		// Empty tiles are white

		// If the laser is on this tile, then we need to draw the laser by filling the tile
		if t & ^NonLaserBits != 0 {
			for i := 0; i < size; i++ {
				for j := 0; j < size; j++ {
					img.Set(x+i, y+j, laserColor)
				}
			}
		}
	case BackslashMirror:
		// Fill the top right if the laser came in on this side
		if t&(LaserLeft|LaserDown) != 0 {
			for i := 0; i < size; i++ {
				for j := 0; j < i; j++ {
					img.Set(x+i, y+j, laserColor)
				}

			}
		}

		// Fill the bottom left if the laser came in on this side
		if t&(LaserRight|LaserUp) != 0 {
			for i := 0; i < size; i++ {
				for j := 0; j < i; j++ {
					img.Set(x+j, y+i, laserColor)
				}

			}
		}

		// Now draw the mirror
		for i := 0; i < size; i++ {
			img.Set(x+i, y+i, color.Black)
		}

	case ForwardSlashMirror:
		// Fill the top left if the laser came in on this side
		if t&(LaserRight|LaserDown) != 0 {
			for i := 0; i < size; i++ {
				for j := 0; j < size-i; j++ {
					img.Set(x+i, y+j, laserColor)
				}

			}
		}

		// Fill the bottom right if the laser came in on this side
		if t&(LaserLeft|LaserUp) != 0 {
			for i := 0; i < size; i++ {
				for j := 0; j < size-i; j++ {
					img.Set(x+size-(i+1), y+size-(j+1), laserColor)
				}

			}
		}

		// Now draw the mirror
		for i := 0; i < size; i++ {
			img.Set(x+i, y+size-(i+1), color.Black)
		}
	case VerticalSplitter:
		middle := x + (size / 2)

		// Fill in the left if the laser came in right, or up / down
		if t&(LaserRight|LaserUp|LaserDown) != 0 {
			for i := x; i < middle; i++ {
				for j := 0; j < size; j++ {
					img.Set(i, y+j, laserColor)
				}
			}
		}

		// Fill in the right if the laser came in left, or up / down
		if t&(LaserLeft|LaserUp|LaserDown) != 0 {
			for i := middle; i < x+size; i++ {
				for j := 0; j < size; j++ {
					img.Set(i, y+j, laserColor)
				}
			}
		}

		// Now draw the mirror
		for i := 0; i < size; i++ {
			img.Set(middle, y+i, color.Black)
		}
	case HorizontalSplitter:
		middle := y + (size / 2)

		// Fill in the top if the laser came in bottom, or left / right
		if t&(LaserDown|LaserLeft|LaserRight) != 0 {
			for i := y; i < middle; i++ {
				for j := 0; j < size; j++ {
					img.Set(x+j, i, laserColor)
				}
			}
		}

		// Fill in the bottom if the laser came in top, or left / right
		if t&(LaserUp|LaserLeft|LaserRight) != 0 {
			for i := middle; i < y+size; i++ {
				for j := 0; j < size; j++ {
					img.Set(x+j, i, laserColor)
				}
			}
		}

		// Now draw the mirror
		for i := 0; i < size; i++ {
			img.Set(x+i, middle, color.Black)
		}
	default:
		panic("unknown tile type")
	}
}
