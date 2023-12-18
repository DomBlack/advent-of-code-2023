package maps

import (
	"fmt"
)

// Direction represents a direction on a map.
type Direction uint8

const (
	// North is the direction up.
	North Direction = iota

	// East is the direction right.
	East

	// South is the direction down.
	South

	// West is the direction left.
	West
)

func (d Direction) FlipHorizontal() Direction {
	switch d {
	case North:
		return North
	case South:
		return South
	case East:
		return West
	case West:
		return East
	default:
		panic(fmt.Sprintf("invalid direction: %d", d))
	}
}
