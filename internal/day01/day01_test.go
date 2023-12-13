package day01

import (
	"testing"
)

func Test_Day1(t *testing.T) {
	input1 := `1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet`

	input2 := `two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen`

	Day01.Test(t, input1, 142, input2, 281)

	Day01.TestPart2(t, "one", 11)
	Day01.TestPart2(t, "two", 22)
	Day01.TestPart2(t, "three", 33)
	Day01.TestPart2(t, "four", 44)
	Day01.TestPart2(t, "five", 55)
	Day01.TestPart2(t, "six", 66)
	Day01.TestPart2(t, "seven", 77)
	Day01.TestPart2(t, "eight", 88)
	Day01.TestPart2(t, "nine", 99)
	Day01.TestPart2(t, "nineight", 98)
}
