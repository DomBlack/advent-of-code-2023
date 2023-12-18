package maps

import (
	"fmt"
)

// Rotate rotates the map in the given direction, which is one of
// [North], [East], [South] & [West], which is
// the equivalent of rotating the map 0, 90, 180 & 270 degrees
//
// Any other direction will result in a panic.
func Rotate[TileType Tile](m *Map[TileType], direction Direction) {
	switch direction {
	case North:
		// no-op, we're already facing up

	case South:
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

	case East:
		for layer := 0; layer < m.Height/2; layer++ {
			first, last := layer, m.Width-1-layer
			for i := first; i < last; i++ {
				offset := i - first

				topIdx := first*m.Width + i
				leftIdx := (last-offset)*m.Width + first
				bottomIdx := last*m.Width + (last - offset)
				rightIdx := i*m.Width + last

				top := m.Tiles[topIdx] // save top
				m.Tiles[topIdx] = m.Tiles[leftIdx]
				m.Tiles[leftIdx] = m.Tiles[bottomIdx]
				m.Tiles[bottomIdx] = m.Tiles[rightIdx]
				m.Tiles[rightIdx] = top // right <- saved top
			}
		}
		m.Width, m.Height = m.Height, m.Width

	case West:
		for layer := 0; layer < m.Height/2; layer++ {
			first, last := layer, m.Width-1-layer
			for i := first; i < last; i++ {
				offset := i - first

				topIdx := first*m.Width + i
				leftIdx := (last-offset)*m.Width + first
				bottomIdx := last*m.Width + (last - offset)
				rightIdx := i*m.Width + last

				top := m.Tiles[topIdx] // save top
				// right -> top
				m.Tiles[topIdx] = m.Tiles[rightIdx]
				// bottom -> right
				m.Tiles[rightIdx] = m.Tiles[bottomIdx]
				// left -> bottom
				m.Tiles[bottomIdx] = m.Tiles[leftIdx]
				// top -> left
				m.Tiles[leftIdx] = top // left <- saved top
			}
		}
		m.Width, m.Height = m.Height, m.Width

	default:
		panic(fmt.Sprintf("unsupported rotation direction: %s", direction))
	}
}
