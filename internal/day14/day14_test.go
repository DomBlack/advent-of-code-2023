package day14

import (
	"testing"
)

func Test_Day14(t *testing.T) {
	input := `
O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....
`

	Day14.Test(t, input, 136, input, 64)
}
