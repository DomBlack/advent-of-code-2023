package maps

import (
	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/vectors/vec2"
)

// Tilt tilts the map in the given direction so that the MovableTile will move in that direction
// until it hits another tile type which isn't an empty tile.
func Tilt[TileType Tile](m *Map[TileType], direction vec2.Vec2, MovableTile TileType) {
	// First rotate the map so we're always facing up
	Rotate(m, direction.FlipHorizontal())

	for idx, tile := range m.Tiles {
		// If we're not looking at a movable tile, then skip it
		if tile != MovableTile {
			continue
		}

		x, y := m.PositionOf(idx)

		// Move the tile towards the top of the map until we hit a non-empty tile
		for y > 0 && m.Tiles[(y-1)*m.Width+x] == m.EmptyType {
			y--
		}

		// Now move the tile to the new position
		newIdx := y*m.Width + x
		if newIdx == idx {
			continue
		}
		m.Tiles[newIdx] = MovableTile
		m.Tiles[idx] = m.EmptyType
	}

	// Finally flip the map back to the original direction
	Rotate(m, direction)
}
