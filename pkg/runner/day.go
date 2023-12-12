package runner

import (
	"fmt"
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
	wrongPart1Answers   []wrongAnswer
	expectedPart1Answer string
	part2               Part[Input]
	wrongPart2Answers   []wrongAnswer
	expectedPart2Answer string
}

// Part represents a function which given the input will return the answer for that part of the day
type Part[Input any] func(log zerolog.Logger, input Input) (answer string, err error)

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
		day:   day,
		part1: part1,
		part2: part2,
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
		day:   day,
		part1: part1,
		part2: part2,
	}
	days[day] = d

	return d
}

func (d *Day[Input, Cache]) Day() int {
	return d.day
}

type wrongAnswer struct {
	value string
	hint  string
}

// WithWrongPart1Answer adds a wrong answer to the list of wrong answers for part 1
func (d *Day[Input, Cache]) WithWrongPart1Answer(answer, hint string) *Day[Input, Cache] {
	d.wrongPart1Answers = append(d.wrongPart1Answers, wrongAnswer{answer, hint})
	return d
}

// WithPart2WrongAnswer adds a wrong answer to the list of wrong answers for part 2
func (d *Day[Input, Cache]) WithPart2WrongAnswer(answer, hint string) *Day[Input, Cache] {
	d.wrongPart2Answers = append(d.wrongPart2Answers, wrongAnswer{answer, hint})
	return d
}

func (d *Day[Input, Cache]) WithExpectedAnswers(part1, part2 string) *Day[Input, Cache] {
	d.expectedPart1Answer = part1
	d.expectedPart2Answer = part2
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

	runPart := func(partNum int, fn Part[Input], expectedAnswer string, wrongAnswers []wrongAnswer) {
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
				if expectedAnswer != "" && answer != expectedAnswer {
					logger.Error().Caller(1).Str("duration", dur.String()).Str("got", answer).Str("expected", expectedAnswer).Msg("part returned wrong answer")
					os.Exit(1)
					return
				}

				for _, wrongAnswer := range wrongAnswers {
					if wrongAnswer.value == answer {
						logger.Warn().Caller(1).Str("duration", dur.String()).Str("answer", answer).Str("hint", wrongAnswer.hint).Msg("part answer incorrect")
						os.Exit(1)
						return
					}
				}

				logger.Info().Caller(1).Str("duration", dur.String()).Str("answer", answer).Msg("part complete")
			}
		} else {
			d.log.Warn().Caller(1).Int("part", 1).Msg("part not implemented")
		}
	}

	// Run part 1 if it exists
	runPart(1, d.part1, d.expectedPart1Answer, d.wrongPart1Answers)
	runPart(2, d.part2, d.expectedPart2Answer, d.wrongPart2Answers)
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
func (d *Day[Input, Cache]) Test(t *testing.T, part1TestInput, part1ExpectedAnswer, part2estInput, part2ExpectedAnswer string) {
	t.Helper()
	t.Parallel()

	d.testPart(t, "part1", 1, d.part1, part1TestInput, part1ExpectedAnswer)
	d.testPart(t, "part2", 2, d.part2, part2estInput, part2ExpectedAnswer)
}

// TestPart1 runs the given part 1 with the given input and asserts the answer
func (d *Day[Input, Cache]) TestPart1(t *testing.T, input, expectedAnswer string) {
	t.Helper()

	d.testPart(t, fmt.Sprintf("part1_%s", input), 1, d.part1, input, expectedAnswer)
}

// TestPart2 runs the given part 2 with the given input and asserts the answer
func (d *Day[Input, Cache]) TestPart2(t *testing.T, input, expectedAnswer string) {
	t.Helper()
	d.testPart(t, fmt.Sprintf("part2_%s", input), 2, d.part2, input, expectedAnswer)
}

// testPart runs the given part with the given input and asserts the answer
func (d *Day[Input, Cache]) testPart(t *testing.T, testName string, partNum int, fn Part[Input], input, expectedAnswer string) {
	t.Helper()
	t.Run(testName, func(t *testing.T) {
		t.Parallel()

		if fn == nil && expectedAnswer == "" {
			t.Skip("Part not implemented")
		}

		if fn == nil {
			assert.FailNow(t, "part not implemented, but an expected answer was provided")
		}
		if expectedAnswer == "" {
			assert.FailNow(t, "part implemented, but not expected answer was provided")
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
