package day03

import (
	"testing"
)

func Test_Temp(t *testing.T) {
	Day03.Run()
}

func Test_Day3(t *testing.T) {
	input := `467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`

	Day03.Test(t, input, "4361", input, "467835")

	Day03.TestPart1(t, `..12..`, "0")
	Day03.TestPart1(t, `*.12..`, "0")
	Day03.TestPart1(t, `.*12..`, "12")
	Day03.TestPart1(t, `..12*.`, "12")
	Day03.TestPart1(t, `..12.*`, "0")

	Day03.TestPart1(t, `
*.....
..12..
......`, "0")

	Day03.TestPart1(t, `
.*....
..12..
......`, "12")

	Day03.TestPart1(t, `
..*...
..12..
......`, "12")

	Day03.TestPart1(t, `
...*..
..12..
......`, "12")

	Day03.TestPart1(t, `
....*.
..12..
......`, "12")

	Day03.TestPart1(t, `
.....*
..12..
......`, "0")
}
