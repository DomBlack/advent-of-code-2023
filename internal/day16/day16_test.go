package day16

import (
	"testing"
)

func Test_Day16(t *testing.T) {
	input := `
.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....
`
	Day16.Test(t, input, 46, input, 51)
}
