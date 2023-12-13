package day09

import (
	"testing"
)

func Test_Day09(t *testing.T) {
	input := `
0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
`

	Day09.Test(t, input, 114, input, 2)
}
