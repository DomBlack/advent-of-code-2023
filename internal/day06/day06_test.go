package day06

import (
	"testing"
)

func Test_Day06(t *testing.T) {
	input := `
Time:      7  15   30
Distance:  9  40  200
`

	Day06.Test(t, input, 288, input, 71503)
}
