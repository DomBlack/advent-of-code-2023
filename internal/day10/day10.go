package day10

import (
	"fmt"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/algorithms/floodfill"
	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day10 = runner.NewDay(10, buildPipeMaze, part1, part2).
	WithExpectedAnswers(6890, 453)

func part1(log zerolog.Logger, input Maze) (answer int, err error) {
	return input.Length / 2, nil
}

func part2(log zerolog.Logger, input Maze) (answer int, err error) {
	input.DebugOutputFile = runner.Output(10)

	enclosedArea, err := input.EnclosedArea()
	if err != nil {
		return 0, errors.Wrap(err, "failed to calculate enclosed area")
	}

	return enclosedArea, nil
}

type Maze struct {
	Start  *Tile
	Size   int
	Length int

	DebugOutputFile string
}

func (m Maze) String() string {
	// Temporarily store our tiles in a 2D array
	tiles := make([][]*Tile, m.Size)
	for i := range tiles {
		tiles[i] = make([]*Tile, m.Size)
	}

	// Walk our maze and populate the 2D array
	prev := m.Start
	tiles[m.Start.X][m.Start.Y] = m.Start
	node := m.Start.Next(nil)
	length := 1
	for node != m.Start {
		tiles[node.X][node.Y] = node

		// These two are done when we make the maze, just here to double check
		if node == nil {
			panic("maze is not a loop")
		}
		if length > 10000 {
			panic("maze is too long")
		}

		node, prev = node.Next(prev), node
		length++
	}

	// Now convert the 2D array into a string
	var sb strings.Builder
	for y := 0; y < m.Size; y++ {
		for x := 0; x < m.Size; x++ {
			if tiles[x][y] == m.Start {
				sb.WriteString("S")
			} else {
				sb.WriteString(tiles[x][y].String())
			}
		}
		sb.WriteString("\n")
	}
	sb.WriteString(fmt.Sprintf("Length: %d\n", m.Length))

	return sb.String()
}

func (m Maze) EnclosedArea() (int, error) {
	// Temporarily store our tiles in a 2D array
	tiles := make([][]*Tile, m.Size)
	for i := range tiles {
		tiles[i] = make([]*Tile, m.Size)
	}

	// Walk our maze and populate the 2D array
	prev := m.Start
	tiles[m.Start.X][m.Start.Y] = m.Start
	node := m.Start.Next(nil)
	length := 1
	for node != m.Start {
		tiles[node.X][node.Y] = node

		// These two are done when we make the maze, just here to double check
		if node == nil {
			panic("maze is not a loop")
		}
		if length > 100_000 {
			panic("maze is too long")
		}

		node, prev = node.Next(prev), node
		length++
	}

	// We're going to construct a 3 by 3 grid of tiles around the start
	//
	// This is so we can fill through the
	const mapScalar = 3

	ff := floodfill.NewMap(m.Size*mapScalar, m.Size*mapScalar)
	for y := 0; y < m.Size; y++ {
		for x := 0; x < m.Size; x++ {
			if tiles[x][y] != nil {
				tile := tiles[x][y]

				switch {
				case tile.North != nil && tile.South != nil:
					pipe := []floodfill.Tile{floodfill.Empty, floodfill.Wall, floodfill.Empty}
					copy(ff.Map[y*mapScalar][x*mapScalar:], pipe)
					copy(ff.Map[y*mapScalar+1][x*mapScalar:], pipe)
					copy(ff.Map[y*mapScalar+2][x*mapScalar:], pipe)

				case tile.East != nil && tile.West != nil:
					empty := []floodfill.Tile{floodfill.Empty, floodfill.Empty, floodfill.Empty}
					copy(ff.Map[y*mapScalar][x*mapScalar:], empty)
					copy(ff.Map[y*mapScalar+1][x*mapScalar:], []floodfill.Tile{floodfill.Wall, floodfill.Wall, floodfill.Wall})
					copy(ff.Map[y*mapScalar+2][x*mapScalar:], empty)

				case tile.North != nil && tile.East != nil:
					copy(ff.Map[y*mapScalar][x*mapScalar:], []floodfill.Tile{floodfill.Empty, floodfill.Wall, floodfill.Empty})
					copy(ff.Map[y*mapScalar+1][x*mapScalar:], []floodfill.Tile{floodfill.Empty, floodfill.Wall, floodfill.Wall})
					copy(ff.Map[y*mapScalar+2][x*mapScalar:], []floodfill.Tile{floodfill.Empty, floodfill.Empty, floodfill.Empty})

				case tile.North != nil && tile.West != nil:
					copy(ff.Map[y*mapScalar][x*mapScalar:], []floodfill.Tile{floodfill.Empty, floodfill.Wall, floodfill.Empty})
					copy(ff.Map[y*mapScalar+1][x*mapScalar:], []floodfill.Tile{floodfill.Wall, floodfill.Wall, floodfill.Empty})
					copy(ff.Map[y*mapScalar+2][x*mapScalar:], []floodfill.Tile{floodfill.Empty, floodfill.Empty, floodfill.Empty})

				case tile.South != nil && tile.West != nil:
					copy(ff.Map[y*mapScalar][x*mapScalar:], []floodfill.Tile{floodfill.Empty, floodfill.Empty, floodfill.Empty})
					copy(ff.Map[y*mapScalar+1][x*mapScalar:], []floodfill.Tile{floodfill.Wall, floodfill.Wall, floodfill.Empty})
					copy(ff.Map[y*mapScalar+2][x*mapScalar:], []floodfill.Tile{floodfill.Empty, floodfill.Wall, floodfill.Empty})

				case tile.South != nil && tile.East != nil:
					copy(ff.Map[y*mapScalar][x*mapScalar:], []floodfill.Tile{floodfill.Empty, floodfill.Empty, floodfill.Empty})
					copy(ff.Map[y*mapScalar+1][x*mapScalar:], []floodfill.Tile{floodfill.Empty, floodfill.Wall, floodfill.Wall})
					copy(ff.Map[y*mapScalar+2][x*mapScalar:], []floodfill.Tile{floodfill.Empty, floodfill.Wall, floodfill.Empty})
				}
			}
		}
	}

	if m.DebugOutputFile != "" {
		ff.ImageScale = 1
		ff.ImageSkipFrames = 9 * 3 * 10 * 2
		ff.StartCapturingFills()
	}

	// Fill from all the edges
	for coord := 0; coord < m.Size*mapScalar; coord++ {
		ff.Fill(coord, 0)
		ff.Fill(0, coord)
		ff.Fill(coord, (m.Size*mapScalar)-1)
		ff.Fill((m.Size*mapScalar)-1, coord)
	}

	if m.DebugOutputFile != "" {
		if err := ff.SaveFillImage(m.DebugOutputFile + ".gif"); err != nil {
			return 0, errors.Wrap(err, "failed to save fill image")
		}
	}

	// Now count the number of tiles which are still marked as empty (i.e. can't be reached from the outside edge of the maze)
	enclosuedTileCount := 0
	for y := 0; y < m.Size; y++ {
		for x := 0; x < m.Size; x++ {
			// nil tiles means this tile wasn't originally a wall of some kind
			if tiles[x][y] == nil {
				if ff.Map[y*mapScalar][x*mapScalar] == floodfill.Empty {
					enclosuedTileCount++
				}
			}
		}
	}

	return enclosuedTileCount, nil
}

type Tile struct {
	X, Y  int
	North *Tile
	East  *Tile
	South *Tile
	West  *Tile
}

// Next returns the next tile in the maze, or nil if there is no next tile
// The prev tile is used to prevent backtracking, if no prev tile is provided
// then a direction will be chosen at random
func (t *Tile) Next(prev *Tile) *Tile {
	if t == nil {
		return nil
	}

	if t.North != nil && t.North != prev {
		return t.North
	} else if t.East != nil && t.East != prev {
		return t.East
	} else if t.South != nil && t.South != prev {
		return t.South
	} else if t.West != nil && t.West != prev {
		return t.West
	}

	return nil
}

func (t *Tile) String() string {
	if t == nil {
		return "."
	}

	switch {
	case t.North != nil && t.South != nil:
		return "|"
	case t.East != nil && t.West != nil:
		return "-"
	case t.North != nil && t.East != nil:
		return "L"
	case t.North != nil && t.West != nil:
		return "J"
	case t.South != nil && t.West != nil:
		return "7"
	case t.South != nil && t.East != nil:
		return "F"
	default:
		return "."
	}
}

func buildPipeMaze(input []byte) (Maze, error) {
	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	if len(lines) == 0 {
		return Maze{}, errors.New("no input")
	}

	size := len(lines[0])
	if size == 0 {
		return Maze{}, errors.New("no line contexts")
	}

	// Create a base 2D array of tiles to work with while parsing
	baseTiles := make([][]*Tile, size)
	for i := range baseTiles {
		baseTiles[i] = make([]*Tile, size)
	}
	getTile := func(x, y int) *Tile {
		if x < 0 || x >= size || y < 0 || y >= size {
			return &Tile{} // temp tile
		}

		tile := baseTiles[x][y]
		if tile == nil {
			tile = &Tile{X: x, Y: y}
			baseTiles[x][y] = tile
		}
		return tile
	}

	var startingTile *Tile

	// Parse our tiles
	for y, line := range lines {
		for x, char := range line {
			tile := getTile(x, y)

			switch char {
			case '|':
				tile.North = getTile(x, y-1)
				tile.South = getTile(x, y+1)
			case '-':
				tile.East = getTile(x+1, y)
				tile.West = getTile(x-1, y)
			case 'L':
				tile.North = getTile(x, y-1)
				tile.East = getTile(x+1, y)
			case 'J':
				tile.North = getTile(x, y-1)
				tile.West = getTile(x-1, y)
			case '7':
				tile.South = getTile(x, y+1)
				tile.West = getTile(x-1, y)
			case 'F':
				tile.South = getTile(x, y+1)
				tile.East = getTile(x+1, y)
			case '.':
				// Do nothing
			case 'S':
				if startingTile != nil {
					return Maze{}, errors.New("multiple starting tiles")
				}
				startingTile = tile
			default:
				return Maze{}, errors.Errorf("unknown tile type: %c", char)
			}
		}
	}

	// Sort out the starting tile
	if startingTile == nil {
		return Maze{}, errors.New("no starting tile found")
	}
	startNorth := getTile(startingTile.X, startingTile.Y-1)
	startEast := getTile(startingTile.X+1, startingTile.Y)
	startSouth := getTile(startingTile.X, startingTile.Y+1)
	startWest := getTile(startingTile.X-1, startingTile.Y)

	if startNorth.South == startingTile {
		startingTile.North = startNorth
	}
	if startEast.West == startingTile {
		startingTile.East = startEast
	}
	if startSouth.North == startingTile {
		startingTile.South = startSouth
	}
	if startWest.East == startingTile {
		startingTile.West = startWest
	}

	// Walk and validate the maze and count the length
	prev := startingTile
	node := startingTile.Next(nil)
	length := 1
	for node != startingTile {
		if node == nil {
			return Maze{}, errors.Newf("maze is not a loop after %d steps", length)
		}

		if length > 100_000 {
			return Maze{}, errors.New("maze is too long")
		}

		node, prev = node.Next(prev), node
		length++
	}

	return Maze{
		Start:  startingTile,
		Size:   size,
		Length: length,
	}, nil
}
