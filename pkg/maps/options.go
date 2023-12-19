package maps

import (
	"image/color"
	"image/draw"
	"math"
)

// mapCfg is the configuration for a map
type mapCfg struct {
	tileRender  TileRender
	tilePalette []color.Color
	minTileSize int
	maxTileSize int
}

// newMapCfg creates a new mapCfg with the default values
func newMapCfg() *mapCfg {
	return &mapCfg{
		tileRender:  fillTileDrawer,
		minTileSize: 1,
		maxTileSize: math.MaxInt,
	}
}

// fillTileDrawer is the default tile drawer which just fills the entire tile with the tile's colour
func fillTileDrawer(tile AnyTile, img draw.Image, x, y, size int) {
	colour := tile.Colour()

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			img.Set(x+i, y+j, colour)
		}
	}
}

// MapOption represents an option that can be applied to a map
type MapOption func(cfg *mapCfg)

// TileRender sets the tile render function for the map
type TileRender func(tile AnyTile, img draw.Image, x, y, size int)

// WithTileRender sets the tile render function for the map
//
// If not set, the map will render each tile fully in the colour of the tile
func WithTileRender(render TileRender) MapOption {
	return func(cfg *mapCfg) {
		cfg.tileRender = render
	}
}

// WithColourPalette sets the colour palette for the map
//
// If not set, the map will compute the palette from the tiles
// by iterating from the zero value, until it gets false from a call to [Tile.Valid]
func WithColourPalette(palette []color.Color) MapOption {
	return func(cfg *mapCfg) {
		cfg.tilePalette = palette
	}
}

// WithMinTileSize sets the minimum tile size for the map  when rendering
func WithMinTileSize(sizeInPixels int) MapOption {
	return func(cfg *mapCfg) {
		cfg.minTileSize = sizeInPixels
	}
}

// WithMaxTileSize sets the maximum tile size for the map when rendering
func WithMaxTileSize(sizeInPixels int) MapOption {
	return func(cfg *mapCfg) {
		cfg.maxTileSize = sizeInPixels
	}
}
