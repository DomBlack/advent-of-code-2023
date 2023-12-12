package floodfill

import (
	"image"
	"image/color"
	"image/gif"
	"os"

	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/ringbuffer"
	"github.com/cockroachdb/errors"
)

var palette = []color.Color{
	color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 0 (Empty): White
	color.RGBA{R: 0, G: 0, B: 0, A: 255},       // 1 (Wall): Black
	color.RGBA{R: 0, G: 255, B: 0, A: 255},     // 2 (Fill): Green
}

// Map represents a map of tiles with a given height and width
type Map struct {
	Height int
	Width  int

	// The map's tiles as a 2D array
	// Map[y][x]
	Map [][]Tile

	ImageScale      int
	ImageSkipFrames int
	fillImages      []*image.Paletted
	fillImageDelays []int
}

// Tile represents a tile in the map
type Tile int

const (
	Empty Tile = iota
	Wall
	Fill
)

// NewMap creates a new map with the given height and width
func NewMap(height, width int) *Map {
	m := &Map{
		Height:          height,
		Width:           width,
		Map:             make([][]Tile, height),
		ImageScale:      1,
		ImageSkipFrames: 27,
	}

	for y := 0; y < height; y++ {
		m.Map[y] = make([]Tile, width)
	}

	return m
}

// Fill fills the map from the given x, y position (assuming it's Empty)
//
// If it is not already empty, nothing will be done
func (m *Map) Fill(x, y int) {
	if !m.isEmpty(x, y) {
		return
	}

	// This uses Span filling from https://en.wikipedia.org/wiki/Flood_fill
	type scanLine struct{ x1, x2, y, dy int }
	queue := ringbuffer.NewGrowable[scanLine]()

	queue.Push(scanLine{x, x, y, 1})
	queue.Push(scanLine{x, x, y - 1, -1})

	count := 0
	set := func(x, y int) {
		m.Map[y][x] = Fill
		count++

		if len(m.fillImages) > 0 && count%m.ImageSkipFrames == 0 {
			m.fillImages = append(m.fillImages, m.Img())
			m.fillImageDelays = append(m.fillImageDelays, 0)
		}
	}

	for {
		line, ok := queue.Dequeue()
		if !ok {
			break
		}
		x := line.x1
		if m.isEmpty(x, line.y) {
			for ; m.isEmpty(x-1, line.y); x-- {
				set(x-1, line.y)
			}

			if x < line.x1 {
				queue.Push(scanLine{x, line.x1 - 1, line.y - line.dy, -line.dy})
			}
		}

		for line.x1 <= line.x2 {
			for ; m.isEmpty(line.x1, line.y); line.x1++ {
				set(line.x1, line.y)
			}

			if line.x1 > x {
				queue.Push(scanLine{x, line.x1 - 1, line.y + line.dy, line.dy})
			}

			if line.x1-1 > line.x2 {
				queue.Push(scanLine{line.x2 + 1, line.x1 - 1, line.y - line.dy, -line.dy})
			}

			line.x1++
			for ; !m.isEmpty(line.x1, line.y) && line.x1 <= line.x2; line.x1++ {
			}
			x = line.x1
		}
	}
}

func (m *Map) isEmpty(x, y int) bool {
	return x >= 0 && x < m.Width && y >= 0 && y < m.Height && m.Map[y][x] == Empty
}

// CountOf returns the number of tiles which are the given type
func (m *Map) CountOf(tileType Tile) int {
	count := 0
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.Map[y][x] == tileType {
				count++
			}
		}
	}

	return count
}

// StartCapturingFills starts capturing the map's during fills
func (m *Map) StartCapturingFills() {
	m.fillImages = []*image.Paletted{m.Img()}
	m.fillImageDelays = []int{10}
}

// SaveFillImage saves the current fill image to the given file
func (m *Map) SaveFillImage(fileName string) error {
	m.fillImages = append(m.fillImages, m.Img())
	m.fillImageDelays = append(m.fillImageDelays, 50)

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to open file to save fill image")
	}
	defer func() { _ = f.Close() }()

	err = gif.EncodeAll(f, &gif.GIF{
		Image:     m.fillImages,
		Delay:     m.fillImageDelays,
		LoopCount: 0,
		Config: image.Config{
			ColorModel: color.Palette(palette),
			Width:      m.Width * m.ImageScale,
			Height:     m.Height * m.ImageScale,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to encode fill image")
	}

	return nil
}

// Img returns an image of the map in it's current state
func (m *Map) Img() *image.Paletted {
	img := image.NewPaletted(image.Rect(0, 0, m.Width*m.ImageScale, m.Height*m.ImageScale), palette)

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			c := palette[m.Map[y][x]]

			for extraX := 0; extraX < m.ImageScale; extraX++ {
				for extraY := 0; extraY < m.ImageScale; extraY++ {
					img.Set(x*m.ImageScale+extraX, y*m.ImageScale+extraY, c)
				}
			}
		}
	}

	return img
}
