package vec2

import (
	"fmt"
)

type Vec2 [2]int

var (
	North = Vec2{0, -1}
	South = Vec2{0, 1}
	West  = Vec2{-1, 0}
	East  = Vec2{1, 0}
)

func (v Vec2) String() string {
	return fmt.Sprintf("(%d, %d)", v[0], v[1])
}

func (v Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v[0] - v2[0], v[1] - v2[1]}
}

// Neg returns the negative of the vector
func (v Vec2) Neg() Vec2 {
	return Vec2{-v[0], -v[1]}
}

func (v Vec2) FlipHorizontal() Vec2 {
	return Vec2{-v[0], v[1]}
}

func (v Vec2) FlipVertical() Vec2 {
	return Vec2{v[0], -v[1]}
}

func (v Vec2) Length() int {
	x := v[0]
	if x < 0 {
		x = -x
	}
	y := v[1]
	if y < 0 {
		y = -y
	}
	return x + y
}
