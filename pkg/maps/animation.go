package maps

import (
	"image"
	"image/color"
	"image/gif"
	"os"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
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
func (m *Map[TileType]) StartCapturingFrames(ctx *runner.Context) {
	if !ctx.SaveOutput() {
		return
	}

	m.captureFrames = true
	m.CaptureFrame("Starting State", 100)
}

// StopCapturingFrames stops capturing frames for the map
func (m *Map[TileType]) StopCapturingFrames(label string) {
	if !m.captureFrames {
		return
	}

	if label == "" {
		label = "Finished"
	}
	m.CaptureFrame(label, 300)

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

func (m *Map[TileType]) SaveAnimationGIF(ctx *runner.Context) error {
	if !ctx.SaveOutput() {
		return nil
	}

	if len(m.Frames) == 0 {
		return errors.New("cannot save animation gif with no frames")
	}

	// End the capture
	if m.captureFrames {
		m.StopCapturingFrames("")
	}

	// Create the output file
	f, err := os.OpenFile(ctx.OutputFile("gif"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to create output file")
	}
	defer func() { _ = f.Close() }()

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
	totalTime := 0
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
		totalTime += frame.delay
	}

	// Set the scale so our animation is at least 500px wide or tall
	scale := 1
	if max(maxWidth, maxHeight) < 500 {
		scale = 500 / max(maxWidth, maxHeight)
	}
	scale = min(max(scale, m.MinTileSize), m.MaxTileSize)

	// Limit our animations to 10 seconds, if they are longer, we want to skip frames
	frameSkip := 1
	if totalTime > 100 {
		frameSkip = int(totalTime / 100)
	}

	// Create the GIF
	imageWidth := maxWidth * scale
	imageHeight := maxHeight * scale
	imageHeight += maxTextHeight
	imageWidth = max(imageWidth, maxTextWidth)
	cfg := &gif.GIF{
		Config: image.Config{
			ColorModel: color.Palette(m.TilePalette),
			Width:      imageWidth,
			Height:     imageHeight,
		},
	}
	for frameIdx, frame := range m.Frames {
		// Skip frames if we need to to stay within our 10 second limit
		// (except for the last frame)
		if frameIdx%frameSkip != 0 && frameIdx != len(m.Frames)-1 {
			continue
		}

		img := image.NewPaletted(image.Rect(0, 0, imageWidth, imageHeight), m.TilePalette)
		for i, tile := range frame.tiles {
			pos := m.PositionOf(i)

			pos[0] *= scale
			pos[1] *= scale

			pos[1] += maxTextHeight

			m.TileRender(tile, img, pos[0], pos[1], scale)
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
