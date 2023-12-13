package day07

import (
	"testing"
)

func Test_Day07(t *testing.T) {
	input := `
32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483
`

	Day07.Test(t, input, 6440, input, 5905)
}
