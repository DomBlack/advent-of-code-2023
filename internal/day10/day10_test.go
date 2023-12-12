package day10

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Day10(t *testing.T) {
	input1 := `
7-F7-
.FJ|7
SJLL7
|F--J
LJ.LJ
`

	input2 := `
.F----7F7F7F7F-7....
.|F--7||||||||FJ....
.||.FJ||||||||L7....
FJL7L7LJLJ||LJ.L-7..
L--J.L7...LJS7F-7L7.
....F-J..F7FJ|L7L7L7
....L7.F7||L7|.L7L7|
.....|FJLJ|FJ|F7|.LJ
....FJL-7.||.||||...
....L---J.LJ.LJLJ...
`

	Day10.Test(t, input1, "8", input2, "")
}

func Test_Day10_Part2(t *testing.T) {
	t.Parallel()

	test := func(name string, expected int, input string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			maze := mustBuildMaze(t, input)
			actual, err := maze.EnclosedArea()
			assert.NoError(t, err, "Error calculating total enclosed area")
			assert.Equal(t, expected, actual, "Total enclosed area")
		})
	}

	test("Example 1", 4, `
 ...........
.S-------7.
.|F-----7|.
.||.....||.
.||.....||.
.|L-7.F-J|.
.|..|.|..|.
.L--J.L--J.
...........
`)

	test("Example 2", 4, `
..........
.S------7.
.|F----7|.
.||....||.
.||....||.
.|L-7F-J|.
.|..||..|.
.L--JL--J.
..........
`)

	test("Example 3", 8, `
.F----7F7F7F7F-7....
.|F--7||||||||FJ....
.||.FJ||||||||L7....
FJL7L7LJLJ||LJ.L-7..
L--J.L7...LJS7F-7L7.
....F-J..F7FJ|L7L7L7
....L7.F7||L7|.L7L7|
.....|FJLJ|FJ|F7|.LJ
....FJL-7.||.||||...
....L---J.LJ.LJLJ...
`)

}

func mustBuildMaze(t *testing.T, input string) Maze {
	maze, err := buildPipeMaze([]byte(input))
	assert.NoError(t, err, "Error building maze")

	_, filename, _, _ := runtime.Caller(1)
	dir := filepath.Dir(filename)

	testdata := filepath.Join(dir, "testdata", t.Name())
	assert.NoError(t, os.MkdirAll(filepath.Dir(testdata), 0755), "Error creating testdata directory")

	maze.DebugOutputFile = testdata

	return maze
}
