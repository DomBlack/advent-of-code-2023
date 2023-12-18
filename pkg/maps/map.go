package maps

import (
	"strings"

	"github.com/cockroachdb/errors"
)

// Map represents a two dimensional map of some
// set of Values.
//
// The zero value of [Tile] is expected to be
// "empty" (i.e. in the case of flood filling, it's
// not a wall). If you want to change this override
// the [EmptyType] field.
type Map[TileType Tile] struct {
	Width  int // The width of the map
	Height int // The height of the map

	// Tiles are all the [Tile]'s that make up the map.
	//
	// The tiles are laid out by rows, such that
	// indexes 0 -> Width are row 1, and indexes
	// Width -> Width * 2 are row 2.
	Tiles []TileType

	// EmptyType represents empty space in the map and
	// is by default the zero value of [Tile].
	EmptyType TileType
}

// From2DSlices creates a new map from the given 2D slice of tiles.
// where the first index is the y position, and the second index is the x position.
func From2DSlices[TileType Tile](tiles [][]TileType) (*Map[TileType], error) {
	if len(tiles) == 0 {
		return nil, errors.New("cannot create map from empty slice")
	}

	height := len(tiles)
	width := len(tiles[0])

	m := &Map[TileType]{
		Height: height,
		Width:  width,
		Tiles:  make([]TileType, height*width),
	}

	// Copy the tiles into the map
	for y, row := range tiles {
		if len(row) != width {
			return nil, errors.Newf("cannot create map from inconsistent row lengths, expected %d tiles per row, got %d", width, len(row))
		}

		for x, tile := range row {
			m.Tiles[y*width+x] = tile
		}
	}

	return m, nil
}

// Tile represents a single tile in a [Map].
type Tile interface {
	~int | ~uint | ~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64 | ~float32 | ~float64

	// Rune should return the single rune that represents
	// this tile.
	Rune() rune
}

// PositionOf returns the x, y position of the given index.
func (m *Map[TileType]) PositionOf(idx int) (x, y int) {
	return idx % m.Width, idx / m.Width
}

// IndexOf returns the index of the given x, y position.
func (m *Map[TileType]) IndexOf(x, y int) int {
	return y*m.Width + x
}

// Get returns the tile at the given x, y position.
//
// If the position is out of bounds, then valid will be false,
// and rtn will be the zero value of [Tile].
//
// Otherwise valid will be true, and rtn will be the tile at
// the given position.
func (m *Map[TileType]) Get(x, y int) (rtn TileType, valid bool) {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return rtn, false
	}

	return m.Tiles[y*m.Width+x], true
}

// Set sets the tile at the given x, y position.
//
// If the position is out of bounds, then valid will be false,
// otherwise valid will be true.
func (m *Map[TileType]) Set(x, y int, tile TileType) (valid bool) {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return false
	}

	m.Tiles[y*m.Width+x] = tile
	return true
}

// String returns a string representation of the map.
func (m *Map[TileType]) String() string {
	var rtn strings.Builder

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			rtn.WriteRune(m.Tiles[y*m.Width+x].Rune())
		}
		rtn.WriteRune('\n')
	}

	return rtn.String()
}
