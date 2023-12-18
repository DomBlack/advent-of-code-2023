package maps

// Tilt tilts the map in the given direction so that the MovableTile will move in that direction
// until it hits another tile type which isn't an empty tile.
func Tilt[TileType Tile](m *Map[TileType], direction Direction, MovableTile TileType) {
	switch direction {
	case North:
		for idx, tile := range m.Tiles {
			// If we're not looking at a movable tile, then skip it
			if tile != MovableTile {
				continue
			}

			y := idx / m.Width

			// Move the tile towards the top of the map until we hit a non-empty tile
			currentIdx := idx
			for y > 0 && m.Tiles[currentIdx-m.Width] == m.EmptyType {
				y--
				currentIdx -= m.Width
			}

			// Now move the tile to the new position
			m.Tiles[idx] = m.EmptyType
			m.Tiles[currentIdx] = MovableTile
		}

	case South:
		maxY := m.Height - 1

		for idx := len(m.Tiles) - 1; idx >= 0; idx-- {
			tile := m.Tiles[idx]

			// If we're not looking at a movable tile, then skip it
			if tile != MovableTile {
				continue
			}

			y := idx / m.Width

			// Move the tile towards the top of the map until we hit a non-empty tile
			currentIdx := idx
			for y < maxY && m.Tiles[currentIdx+m.Width] == m.EmptyType {
				y++
				currentIdx += m.Width
			}

			// Now move the tile to the new position
			m.Tiles[idx] = m.EmptyType
			m.Tiles[currentIdx] = MovableTile
		}

	case East:
		maxX := m.Width - 1

		for idx := len(m.Tiles) - 1; idx >= 0; idx-- {
			tile := m.Tiles[idx]
			// If we're not looking at a movable tile, then skip it
			if tile != MovableTile {
				continue
			}

			x := idx % m.Width

			// Move the tile towards the top of the map until we hit a non-empty tile
			currentIdx := idx
			for x < maxX && m.Tiles[currentIdx+1] == m.EmptyType {
				x++
				currentIdx++
			}

			// Now move the tile to the new position
			m.Tiles[idx] = m.EmptyType
			m.Tiles[currentIdx] = MovableTile
		}

	case West:
		for idx, tile := range m.Tiles {
			// If we're not looking at a movable tile, then skip it
			if tile != MovableTile {
				continue
			}

			x := idx % m.Width

			// Move the tile towards the top of the map until we hit a non-empty tile
			currentIdx := idx
			for x > 0 && m.Tiles[currentIdx-1] == m.EmptyType {
				x--
				currentIdx--
			}

			// Now move the tile to the new position
			m.Tiles[idx] = m.EmptyType
			m.Tiles[currentIdx] = MovableTile
		}

	default:
		panic("invalid direction")
	}
}
