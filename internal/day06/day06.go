package day06

import (
	"sort"
	"strconv"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day06 = runner.NewDay(6, parseBaseInput, part1, part2)

func part1(log zerolog.Logger, input baseInput) (answer int, err error) {
	result := 1

	times := strings.Fields(input.TimeLine)
	distances := strings.Fields(input.DistanceLine)

	if len(times) != len(distances) {
		return 0, errors.Newf("mismatched number of times and distances: %d != %d", len(times), len(distances))
	}

	for i := range times {
		race, err := parseRace(times[i], distances[i])
		if err != nil {
			return 0, err
		}

		result *= race.NumberWinningMethods()

	}

	return result, nil
}

func part2(log zerolog.Logger, input baseInput) (answer int, err error) {
	race, err := parseRace(
		strings.ReplaceAll(input.TimeLine, " ", ""),
		strings.ReplaceAll(input.DistanceLine, " ", ""),
	)
	if err != nil {
		return 0, err
	}

	return race.NumberWinningMethods(), nil
}

type Race struct {
	Time           int
	RecordDistance int
}

func (r Race) NumberWinningMethods() int {
	// Binary search the first time we win
	firstWin := sort.Search(r.Time+1, func(i int) bool {
		return i*(r.Time-i) > r.RecordDistance
	})

	// We never win!
	if firstWin == r.Time+1 {
		return 0
	}

	// Binary search the first time we lose after winning
	firstLose := sort.Search(r.Time+1-firstWin, func(i int) bool {
		return i*((r.Time-i)+firstWin) < r.RecordDistance
	})

	wins := firstLose - firstWin

	return wins
}

func parseRace(time, recordDistance string) (Race, error) {
	timeInt, err := strconv.Atoi(time)
	if err != nil {
		return Race{}, errors.Wrapf(err, "failed to parse time %q", time)
	}

	recordDistanceInt, err := strconv.Atoi(recordDistance)
	if err != nil {
		return Race{}, errors.Wrapf(err, "failed to parse distance %q", recordDistance)
	}

	return Race{Time: timeInt, RecordDistance: recordDistanceInt}, nil
}

type baseInput struct {
	TimeLine     string
	DistanceLine string
}

func parseBaseInput(input []byte) (baseInput, error) {
	timeLine, distanceLine, found := strings.Cut(strings.TrimSpace(string(input)), "\n")
	if !found {
		return baseInput{}, errors.Newf("no newline found in input: %q", string(input))
	}

	timeLine = strings.TrimSpace(strings.TrimPrefix(timeLine, "Time:"))
	distanceLine = strings.TrimSpace(strings.TrimPrefix(distanceLine, "Distance:"))

	return baseInput{TimeLine: strings.TrimSpace(timeLine), DistanceLine: strings.TrimSpace(distanceLine)}, nil
}
