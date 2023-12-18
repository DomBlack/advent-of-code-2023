package maps

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// frame represents a single frame of a map
type frame[TileType Tile] struct {
	label  string     // the label of the frame
	width  int        // the width of the map as it was on this frame
	height int        // the height of the map as it was on this frame
	delay  int        // the delay in 100ths of a second before the next frame should be displayed
	tiles  []TileType // the tiles of the map
}

// StartCapturingFrames starts capturing frames for the map
// starting with the current state of the map
func (m *Map[TileType]) StartCapturingFrames() {
	m.captureFrames = true
	m.CaptureFrame("Starting State", 100)
}

// StopCapturingFrames stops capturing frames for the map
func (m *Map[TileType]) StopCapturingFrames(label string) {
	if label == "" {
		label = "Finished"
	}
	m.CaptureFrame(label, 100)

	m.captureFrames = false
}

// CaptureFrame captures the current state of the map as a frame
// with the given label and delay if and only if we are currently
// capturing frames. Otherwise this function does nothing.
func (m *Map[TileType]) CaptureFrame(label string, delay int) {
	if !m.captureFrames {
		return
	}

	m.Frames = append(m.Frames, frame[TileType]{
		label:  label,
		width:  m.Width,
		height: m.Height,
		delay:  delay,
		tiles:  append([]TileType(nil), m.Tiles...),
	})
}

func (m *Map[TileType]) SaveAnimationGIF(fileName string) error {
	if len(m.Frames) == 0 {
		return errors.New("cannot save animation gif with no frames")
	}

	// End the capture
	if m.captureFrames {
		m.StopCapturingFrames("")
	}

	// Create the output file
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to create output file")
	}
	defer func() { _ = f.Close() }()

	// Build our palette using the source tiles
	palette := make([]color.Color, 0)
	var tile TileType
	for {
		if !tile.Valid() {
			break
		}

		palette = append(palette, tile.Colour())
		tile++
	}
	// always add black for text
	palette = append(palette, color.Black)

	fontToUse := basicfont.Face7x13
	fontDraw := &font.Drawer{
		Src:  image.NewUniform(color.Black),
		Face: fontToUse,
		Dot:  fixed.Point26_6{X: fixed.I(0), Y: fixed.I(0)},
	}

	// Calculate the max width and height of the frames
	maxWidth := 0
	maxHeight := 0
	maxTextWidth := 0
	maxTextHeight := 0
	var totalTime time.Duration
	for _, frame := range m.Frames {
		if frame.width > maxWidth {
			maxWidth = frame.width
		}
		if frame.height > maxHeight {
			maxHeight = frame.height
		}

		if frame.label != "" {
			bounds := fontDraw.MeasureString(frame.label)

			w := bounds.Ceil()
			if w > maxTextWidth {
				maxTextWidth = w
			}

			maxTextHeight = 18
		}

		// each delay is in 100ths of a second
		totalTime += (10 * time.Millisecond) * time.Duration(frame.delay)
	}

	// Set the scale so our animation is at least 500px wide or tall
	scale := 1
	if max(maxWidth, maxHeight) < 500 {
		scale = 500 / max(maxWidth, maxHeight)
	}

	// Limit our animations to 10 seconds, if they are longer, we want to skip frames
	frameSkip := 1
	if totalTime > 10*time.Second {
		frameSkip = int(totalTime / (10 * time.Second))
	}

	// Create the GIF
	imageWidth := maxWidth * scale
	imageHeight := maxHeight * scale
	imageHeight += maxTextHeight
	imageWidth = max(imageWidth, maxTextWidth)
	cfg := &gif.GIF{
		Config: image.Config{
			ColorModel: color.Palette(palette),
			Width:      imageWidth,
			Height:     imageHeight,
		},
	}
	for frameIdx, frame := range m.Frames {
		if frameIdx%frameSkip != 0 {
			continue
		}

		img := image.NewPaletted(image.Rect(0, 0, imageWidth, imageHeight), palette)
		for i, tile := range frame.tiles {
			x, y := m.PositionOf(i)

			x *= scale
			y *= scale

			y += maxTextHeight

			tileColour := tile.Colour()
			for scaleX := 0; scaleX < scale; scaleX++ {
				for scaleY := 0; scaleY < scale; scaleY++ {
					img.Set(x+scaleX, y+scaleY, tileColour)
				}
			}
		}

		if frame.label != "" {
			d := &font.Drawer{
				Dst:  img,
				Src:  image.NewUniform(color.Black),
				Face: fontToUse,
				Dot:  fixed.Point26_6{X: fixed.I(5), Y: fixed.I(15)},
			}
			d.DrawString(frame.label)
		}

		cfg.Image = append(cfg.Image, img)
		cfg.Delay = append(cfg.Delay, frame.delay)
	}

	// Save the GIF
	if err := gif.EncodeAll(f, cfg); err != nil {
		return errors.Wrap(err, "failed to encode gif")
	}

	// Clear the frames
	m.Frames = nil

	return nil
}
