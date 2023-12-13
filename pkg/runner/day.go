package runner

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

// Day represents a day of the advent of code
type Day[Input, Cache any] struct {
	log                 zerolog.Logger
	day                 int
	inputPreprocessor   func([]byte) (Cache, error)
	cacheToInput        func(Cache) Input
	part1               Part[Input]
	part1Answers        answers
	part2               Part[Input]
	part2Answers        answers
	expectedPart2Answer *int
}

type answers struct {
	min, max int
	answer   *int
}

// Part represents a function which given the input will return the answer for that part of the day
type Part[Input any] func(log zerolog.Logger, input Input) (answer int, err error)

// NewStreamingDay creates a day which takes a parser to convert the input bytes into a stream of inputs,
// which the part 1 and part 2 functions can then use.
func NewStreamingDay[Input any](day int, parser func([]byte) stream.Stream[Input], part1, part2 Part[stream.Stream[Input]]) *Day[stream.Stream[Input], []Input] {
	d := &Day[stream.Stream[Input], []Input]{
		log: log.With().Int("_day", day).Logger().Level(zerolog.InfoLevel),
		inputPreprocessor: func(bytes []byte) ([]Input, error) {
			return stream.Collect(parser(bytes))
		},
		cacheToInput: func(cache []Input) stream.Stream[Input] {
			return stream.From(cache)
		},
		day:          day,
		part1:        part1,
		part1Answers: answers{min: math.MinInt, max: math.MaxInt},
		part2:        part2,
		part2Answers: answers{min: math.MinInt, max: math.MaxInt},
	}

	days[day] = d

	return d
}

// NewDay returns a day which is parsed as a whole initially, and then the
// parsed data is given to each part.
//
// See [NewStreamingDay] for a stream processing version
func NewDay[Input any](day int, parser func([]byte) (Input, error), part1, part2 Part[Input]) *Day[Input, Input] {
	d := &Day[Input, Input]{
		log:               log.With().Int("_day", day).Logger().Level(zerolog.InfoLevel),
		inputPreprocessor: parser,
		cacheToInput: func(cache Input) Input {
			return cache
		},
		day:          day,
		part1:        part1,
		part1Answers: answers{min: math.MinInt, max: math.MaxInt},
		part2:        part2,
		part2Answers: answers{min: math.MinInt, max: math.MaxInt},
	}
	days[day] = d

	return d
}

func (d *Day[Input, Cache]) Day() int {
	return d.day
}

func (d *Day[Input, Cache]) WithExpectedAnswers(part1, part2 int) *Day[Input, Cache] {
	d.part1Answers.answer = &part1
	d.part2Answers.answer = &part2
	return d
}

func (d *Day[Input, Cache]) DisableLogging() {
	d.log = d.log.Level(zerolog.Disabled)
}

// Run executes the given parts with the given input
//
// If an error is encountered, the program will exit with a non-zero exit code
func (d *Day[Input, Cache]) Run() {
	// Read the input
	readStart := time.Now()
	inputFile := filepath.Join(repoDir, "inputs", fmt.Sprintf("day%02d.txt", d.day))
	input, err := os.ReadFile(inputFile)
	if err != nil {
		d.log.Err(err).Str("file", inputFile).Msg("failed to read input file")
		os.Exit(1)
		return
	}

	// Preprocess the input
	cacheData, err := d.inputPreprocessor(input)
	if err != nil {
		d.log.Err(err).Msg("failed to preprocess input")
		os.Exit(1)
		return
	}
	d.log.Info().Str("duration", time.Since(readStart).String()).Msg("days input parsed")

	runPart := func(partNum int, fn Part[Input], answers answers) {
		if fn != nil {
			logger := d.log.With().Int("_part", 1).Logger()

			start := time.Now()
			answer, err := fn(logger, d.cacheToInput(cacheData))
			dur := time.Since(start)
			if err != nil {
				logger.Err(err).Str("duration", dur.String()).Msg("failed to run part")
				os.Exit(1)
				return
			} else {
				if answers.answer != nil && answer != *answers.answer {
					logger.Error().Caller(1).Str("duration", dur.String()).Int("got", answer).Int("expected", *answers.answer).Msg("part returned wrong answer")
					os.Exit(1)
					return
				}

				if answers.min >= answer {
					logger.Error().Caller(1).Str("duration", dur.String()).Int("got", answer).Int("min", answers.min).Msg("part returned answer below hinted minimum")
					os.Exit(1)
					return
				}

				if answers.max < answer {
					logger.Error().Caller(1).Str("duration", dur.String()).Int("got", answer).Int("max", answers.max).Msg("part returned answer above hinted maximum")
					os.Exit(1)
					return
				}

				logger.Info().Caller(1).Str("duration", dur.String()).Int("answer", answer).Msg("part complete")
			}
		} else {
			d.log.Warn().Caller(1).Int("part", 1).Msg("part not implemented")
		}
	}

	// Run part 1 if it exists
	runPart(1, d.part1, d.part1Answers)
	runPart(2, d.part2, d.part2Answers)
}

// Output returns an output path for the day
func Output(day int) string {
	dir := filepath.Join(repoDir, "outputs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal().Err(err).Str("dir", dir).Msg("failed to create output directory")
	}

	return filepath.Join(dir, fmt.Sprintf("day%02d", day))
}

// Test runs the given parts with the given input and asserts the answers
func (d *Day[Input, Cache]) Test(t *testing.T, part1TestInput string, part1ExpectedAnswer int, part2estInput string, part2ExpectedAnswer int) {
	t.Helper()
	t.Parallel()

	d.testPart(t, "part1", 1, d.part1, part1TestInput, part1ExpectedAnswer)
	d.testPart(t, "part2", 2, d.part2, part2estInput, part2ExpectedAnswer)
}

// TestPart1 runs the given part 1 with the given input and asserts the answer
func (d *Day[Input, Cache]) TestPart1(t *testing.T, input string, expectedAnswer int) {
	t.Helper()

	d.testPart(t, fmt.Sprintf("part1_%s", input), 1, d.part1, input, expectedAnswer)
}

// TestPart2 runs the given part 2 with the given input and asserts the answer
func (d *Day[Input, Cache]) TestPart2(t *testing.T, input string, expectedAnswer int) {
	t.Helper()
	d.testPart(t, fmt.Sprintf("part2_%s", input), 2, d.part2, input, expectedAnswer)
}

// testPart runs the given part with the given input and asserts the answer
func (d *Day[Input, Cache]) testPart(t *testing.T, testName string, partNum int, fn Part[Input], input string, expectedAnswer int) {
	t.Helper()
	t.Run(testName, func(t *testing.T) {
		t.Parallel()

		if fn == nil {
			t.Skip("Part not implemented")
		}
		if fn == nil && expectedAnswer != 0 {
			assert.FailNow(t, "answer provided, but part not implemented yet")
		}

		// drop to trace level for tests
		testLogger := d.log.Level(zerolog.TraceLevel).With().Int("_part", partNum).Logger()

		preppedData, err := d.inputPreprocessor([]byte(strings.TrimSpace(input)))
		assert.NoError(t, err, "Failed to preprocess input")

		answer, err := fn(testLogger, d.cacheToInput(preppedData))
		assert.NoError(t, err)
		assert.Equal(t, expectedAnswer, answer, "Part answer incorrect")
	})
}
