package polygonarea

import (
	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/vectors/vec2"
)

// Area calculates the area of a polygon with integer coordinates
//
// Note: This function requires that all points are one unit apart
// i.e. you can't go from (0, 0) to (0, 10) in one step, but must go (0, 0) -> (0, 1) -> (0, 2) -> ... -> (0, 10)
func Area(points []vec2.Vec2) (area int) {
	n := len(points)
	if n < 3 {
		return 0
	}

	// Shoelace formula to calculate the area of a polygon
	// https://en.wikipedia.org/wiki/Shoelace_formula
	sum := 0
	for i, current := range points {
		next := points[(i+1)%n]
		sum += (current[1] + next[1]) * (current[0] - next[0])
	}
	if sum < 0 {
		sum = -sum
	}
	area = sum / 2.0

	// Pick's theorem to calculate the area of a polygon when using integer coordinates
	// https://en.wikipedia.org/wiki/Pick%27s_theorem
	return area + len(points)/2 + 1
}
