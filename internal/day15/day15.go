package day15

import (
	"strconv"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day15 = runner.NewStreamingDay(15, parseCommas, part1, part2)

func parseCommas(bytes []byte) stream.Stream[string] {
	return stream.SplitBy(bytes, ',')
}

func part1(log zerolog.Logger, input stream.Stream[string]) (answer int, err error) {
	hashed := stream.Map(input, hash)
	return stream.Sum(hashed)
}

func part2(log zerolog.Logger, input stream.Stream[string]) (answer int, err error) {
	type pair struct {
		label string
		value int
	}
	boxes := [256][]pair{}

	err = stream.ForEach(input, func(input string) error {
		if input[len(input)-1] == '-' {
			label := input[:len(input)-1]
			labelHash, _ := hash(label)

			newList := make([]pair, 0, len(boxes[labelHash]))
			for _, pair := range boxes[labelHash] {
				if pair.label != label {
					newList = append(newList, pair)
				}
			}
			boxes[labelHash] = newList
		} else {
			label, valueStr, found := strings.Cut(input, "=")
			if !found {
				return errors.Newf("Invalid input: %s", input)
			}
			labelHash, _ := hash(label)

			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return errors.Wrapf(err, "Invalid lens value: %s", valueStr)
			}

			if value < 1 || value > 9 {
				return errors.Newf("Invalid lens value: %d", value)
			}

			found = false
			for i, pair := range boxes[labelHash] {
				if pair.label == label {
					found = true
					boxes[labelHash][i].value = value
					break
				}
			}

			if !found {
				boxes[labelHash] = append(boxes[labelHash], pair{
					label: label,
					value: value,
				})
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	for i, box := range boxes {
		for slot, pair := range box {
			power := (i + 1) * (slot + 1) * pair.value
			log.Debug().Msgf("%s: %d (box %d) * %d (%d slot) * %d (focal length) = %d", pair.label, i+1, i, slot+1, slot, pair.value, power)
			answer += power
		}
	}

	return answer, nil
}

func hash(input string) (rtn int, err error) {
	for _, char := range input {
		rtn += int(char)
		rtn *= 17
		rtn %= 256
	}
	return rtn, nil
}
