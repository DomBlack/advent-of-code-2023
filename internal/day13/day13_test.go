package day13

import (
	"testing"
)

func Test_Day13(t *testing.T) {
	input := `
#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
`

	Day13.Test(t, input, 405, input, 400)
}
