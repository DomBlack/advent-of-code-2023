package day12

import (
	"fmt"
	"testing"

	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/stretchr/testify/assert"
)

func Test_Day12(t *testing.T) {
	input := `
???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1
`
	Day12.Test(t, input, "21", input, "525152")
}

func Test_Day12_Options(t *testing.T) {
	t.Parallel()

	test := func(input string, unfold bool, expected int) {
		t.Helper()

		t.Run(input, func(t *testing.T) {
			t.Parallel()

			springs, err := stream.Collect(parseInput([]byte(input)))
			assert.NoError(t, err, "Failed to parse input")
			assert.Len(t, springs, 1, "Expected only one spring")

			spring := springs[0]
			if unfold {
				spring = spring.Unfold()
				fmt.Println(spring.String())
			}

			counts, err := countOptions(make(map[string]int), spring.Springs, spring.DamagedSpringGroups)
			assert.NoError(t, err, "Failed to count options")
			assert.Equal(t, counts, expected, "Expected %d options, got %d", expected, counts)
		})
	}

	test("???.### 1,1,3", false, 1)
	test(".??..??...?##. 1,1,3", false, 4)
	test("?#?#?#?#?#?#?#? 1,3,1,6", false, 1)
	test("????.#...#... 4,1,1", false, 1)
	test("????.######..#####. 1,6,5", false, 4)
	test("?###???????? 3,2,1", false, 10)

	test(".# 1", true, 1)
	test("???.### 1,1,3", true, 1)
	test(".??..??...?##. 1,1,3", true, 16384)
	test("?#?#?#?#?#?#?#? 1,3,1,6", true, 1)
	test("????.#...#... 4,1,1", true, 16)
	test("????.######..#####. 1,6,5", true, 2500)
	test("?###???????? 3,2,1", true, 506250)
}
