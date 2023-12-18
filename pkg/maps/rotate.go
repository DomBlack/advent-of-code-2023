package maps

import (
	"fmt"

	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/vectors/vec2"
)

// Rotate rotates the map in the given direction, which is one of
// [vec2.North], [vec2.East], [vec2.South] & [vec2.West], which is
// the equivalent of rotating the map 0, 90, 180 & 270 degrees
//
// Any other direction will result in a panic.
func Rotate[TileType Tile](m *Map[TileType], direction vec2.Vec2) {
	switch direction {
	case vec2.North:
		// no-op, we're already facing up

	case vec2.South:
		for y := 0; y < m.Height/2; y++ {
			for x := 0; x < m.Width; x++ {
				originalIdx := y*m.Width + x
				rotatedIdx := (m.Height-y-1)*m.Width + (m.Width - x - 1)
				m.Tiles[originalIdx], m.Tiles[rotatedIdx] = m.Tiles[rotatedIdx], m.Tiles[originalIdx]
			}
		}

		// If the height is odd, then we need to rotate the middle row
		if m.Height%2 == 1 {
			mid := (m.Height / 2) + 1
			for x := 0; x < m.Width/2; x++ {
				originalIdx := mid*m.Width + x
				rotatedIdx := mid*m.Width + (m.Width - x - 1)
				m.Tiles[originalIdx], m.Tiles[rotatedIdx] = m.Tiles[rotatedIdx], m.Tiles[originalIdx]
			}
		}

	case vec2.East:
		newTiles := make([]TileType, m.Width*m.Height)
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				originalIdx := y*m.Width + x
				rotatedIdx := x*m.Height + (m.Height - y - 1)
				newTiles[rotatedIdx] = m.Tiles[originalIdx]
			}
		}

		m.Tiles = newTiles
		m.Width, m.Height = m.Height, m.Width

	case vec2.West:
		newTiles := make([]TileType, m.Width*m.Height)
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				originalIdx := y*m.Width + x
				rotatedIdx := (m.Width-x-1)*m.Height + y
				newTiles[rotatedIdx] = m.Tiles[originalIdx]
			}
		}

		m.Tiles = newTiles
		m.Width, m.Height = m.Height, m.Width
	default:
		panic(fmt.Sprintf("unsupported rotation direction: %s", direction))
	}
}
