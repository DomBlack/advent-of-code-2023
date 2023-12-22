package floodfill

import (
	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/ringbuffer"
	"github.com/DomBlack/advent-of-code-2023/pkg/maps"
)

type scanLine struct{ x1, x2, y, dy int }

// Fill fills the map from the given x, y position (assuming it's Empty)
func Fill[TileType maps.Tile](m *maps.Map[TileType], from maps.Pos, fillTile TileType) {
	var emptyTile TileType

	isEmpty := func(x, y int) bool {
		tile, valid := m.Get(maps.Pos{x, y})
		return valid && tile == emptyTile
	}

	if !isEmpty(from[0], from[1]) {
		return
	}

	// This uses Span filling from https://en.wikipedia.org/wiki/Flood_fill
	queue := ringbuffer.NewGrowable[scanLine]()

	queue.Push(scanLine{from[0], from[0], from[1], 1})
	queue.Push(scanLine{from[0], from[0], from[1] - 1, -1})

	count := 0
	set := func(x, y int) {
		m.Set(maps.Pos{x, y}, fillTile)
		count++
	}

	for {
		line, ok := queue.Dequeue()
		if !ok {
			break
		}
		x := line.x1
		if isEmpty(x, line.y) {
			count = 0
			for ; isEmpty(x-1, line.y); x-- {
				set(x-1, line.y)
			}
			if count > 0 {
				m.CaptureFrame("", 1)
			}

			if x < line.x1 {
				queue.Push(scanLine{x, line.x1 - 1, line.y - line.dy, -line.dy})
			}
		}

		for line.x1 <= line.x2 {
			count = 0
			for ; isEmpty(line.x1, line.y); line.x1++ {
				set(line.x1, line.y)
			}
			if count > 0 {
				m.CaptureFrame("", 1)
			}

			if line.x1 > x {
				queue.Push(scanLine{x, line.x1 - 1, line.y + line.dy, line.dy})
			}

			if line.x1-1 > line.x2 {
				queue.Push(scanLine{line.x2 + 1, line.x1 - 1, line.y - line.dy, -line.dy})
			}

			line.x1++
			for ; !isEmpty(line.x1, line.y) && line.x1 <= line.x2; line.x1++ {
			}
			x = line.x1
		}
	}
}
