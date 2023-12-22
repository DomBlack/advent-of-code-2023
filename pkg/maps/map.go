package maps

import (
	"image/color"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/vectors/vec2"
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

	captureFrames bool              // Are we capturing frames?
	TileRender    TileRender        // How to render the tiles
	TilePalette   []color.Color     // The palette of colours to use when rendering the tiles
	MinTileSize   int               // The minimum size of the map when rendered
	MaxTileSize   int               // The maximum size of the map when rendered
	Frames        []frame[TileType] // The frames of the we've captured
}

// Pos represents a position on the map with x, y coordinates.
type Pos = vec2.Vec2

// New creates a new map with the given width and height will all tiles set to the zero value of [Tile].
func New[TileType Tile](width, height int, options ...MapOption) *Map[TileType] {
	cfg := newMapCfg()
	for _, option := range options {
		option(cfg)
	}

	m := &Map[TileType]{
		Height:      height,
		Width:       width,
		TileRender:  cfg.tileRender,
		TilePalette: cfg.tilePalette,
		MinTileSize: cfg.minTileSize,
		MaxTileSize: cfg.maxTileSize,
		Tiles:       make([]TileType, height*width),
	}

	// If the options didn't include a colour palette, then automatically
	// generate one.
	if len(cfg.tilePalette) == 0 {
		var tile TileType
		for {
			if !tile.Valid() {
				break
			}

			m.TilePalette = append(m.TilePalette, tile.Colour())
			tile++
		}
	}
	// always add black for text
	m.TilePalette = append(m.TilePalette, color.Black)

	return m
}

// From2DSlices creates a new map from the given 2D slice of tiles.
// where the first index is the y position, and the second index is the x position.
func From2DSlices[TileType Tile](tiles [][]TileType, options ...MapOption) (*Map[TileType], error) {
	if len(tiles) == 0 {
		return nil, errors.New("cannot create map from empty slice")
	}

	height := len(tiles)
	width := len(tiles[0])

	m := New[TileType](width, height, options...)

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
	~int | ~uint | ~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64
	AnyTile
}

// AnyTile represents an untyped tile that can be used in a [Map].
type AnyTile interface {
	// Valid returns true if this tile is a valid tile.
	Valid() bool

	// Rune should return the single rune that represents
	// this tile.
	Rune() rune

	// Colour is the colour of the tile when displayed
	// in an animation.
	Colour() color.Color
}

// PositionOf returns the x, y position of the given index.
func (m *Map[TileType]) PositionOf(idx int) Pos { return Pos{idx % m.Width, idx / m.Width} }

// IndexOf returns the index of the given x, y position.
func (m *Map[TileType]) IndexOf(pos Pos) int { return pos[1]*m.Width + pos[0] }

// InBounds returns true if the given position is in bounds of the map.
func (m *Map[TileType]) InBounds(pos Pos) bool {
	return pos[0] >= 0 && pos[0] < m.Width && pos[1] >= 0 && pos[1] < m.Height
}

// Neighbours returns the neighbours of the given position.
func (m *Map[TileType]) Neighbours(pos Pos) (rtn []Pos) {
	for _, offset := range vec2.CardinalOffsets {
		if nPos := pos.Add(offset); m.InBounds(nPos) {
			rtn = append(rtn, nPos)
		}
	}

	return rtn
}

// Get returns the tile at the given x, y position.
//
// If the position is out of bounds, then valid will be false,
// and rtn will be the zero value of [Tile].
//
// Otherwise valid will be true, and rtn will be the tile at
// the given position.
func (m *Map[TileType]) Get(pos Pos) (rtn TileType, valid bool) {
	if !m.InBounds(pos) {
		return m.EmptyType, false
	}

	return m.Tiles[pos[1]*m.Width+pos[0]], true
}

// Set sets the tile at the given x, y position.
//
// If the position is out of bounds, then valid will be false,
// otherwise valid will be true.
func (m *Map[TileType]) Set(pos Pos, tile TileType) (valid bool) {
	if !m.InBounds(pos) {
		return false
	}

	m.Tiles[pos[1]*m.Width+pos[0]] = tile
	return true
}

// AddFlagAt adds the given flag to the tile at the given x, y position.
func (m *Map[TileType]) AddFlagAt(pos Pos, flag TileType) {
	if m.InBounds(pos) {
		m.Tiles[pos[1]*m.Width+pos[0]] |= flag
	}
}

// RemoveFlagAt removes the given flag from the tile at the given x, y position.
func (m *Map[TileType]) RemoveFlagAt(pos Pos, flag TileType) {
	if m.InBounds(pos) {
		m.Tiles[pos[1]*m.Width+pos[0]] &^= flag
	}
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
