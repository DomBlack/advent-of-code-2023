package day03

import (
	"io"
	"strconv"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day03 = runner.NewStreamingDay(3, parseSchematic, part1, part2).
	WithExpectedAnswers(539590, 80703636)

func part1(_ *runner.Context, _ zerolog.Logger, input stream.Stream[Token]) (answer int, err error) {
	parts, numbers := stream.Partition(input, func(token Token) (bool, error) {
		if token.Type == Part {
			return true, nil
		} else if token.Type == Number {
			return false, nil
		} else {
			return false, errors.Newf("unknown token type: %s", token.Type)
		}
	})

	resettableParts := stream.Resettable(parts)
	partNumbers := stream.Filter(numbers, func(number Token) (bool, error) {
		resettableParts.Restore() // Restore to our previous save point

		for {
			part, err := resettableParts.Next()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return false, nil
				}
				return false, err
			}

			lineDiff := part.Line - number.Line
			if lineDiff > 1 {
				// No point in continuing to search for this part number as it's too far away
				return false, nil
			} else if lineDiff < -1 {
				// If the part is more than two lines above the number then we can skip it
				// but we can also save this position to restore for the next number as we
				// know it's too far away for that too
				resettableParts.Save()
			} else {
				numLeft := number.StartCol()
				numRight := number.EndCol()
				partLeft := part.StartCol() - 1
				partRight := part.EndCol() + 1

				if (numLeft <= partLeft && partLeft <= numRight) || (numLeft <= partRight && partRight <= numRight) {
					// log.Debug().Str("part_number", number.String()).Str("part_symbol", part.String()).Int("part_left", partLeft).Int("part_right", partRight).Int("num_left", numLeft).Int("num_right", numRight).Msg("Found overlap")
					return true, nil
				}
			}
		}
	})

	partNumbersAsInts := stream.Map(partNumbers, func(number Token) (int, error) {
		return strconv.Atoi(number.Value)
	})

	return stream.Sum(partNumbersAsInts)
}

func part2(_ *runner.Context, _ zerolog.Logger, input stream.Stream[Token]) (answer int, err error) {
	parts, numbers := stream.Partition(input, func(token Token) (bool, error) {
		if token.Type == Part {
			return true, nil
		} else if token.Type == Number {
			return false, nil
		} else {
			return false, errors.Newf("unknown token type: %s", token.Type)
		}
	})

	possibleGears := stream.Filter(parts, func(part Token) (bool, error) {
		return part.Value == "*", nil
	})
	possibleGearPtrs := stream.Map(possibleGears, func(part Token) (*Token, error) {
		return &part, nil
	})
	resettableGearPtrs := stream.Resettable(possibleGearPtrs)

	// Attach all the numbers to the possible gears which are adjacent to them
	err = stream.ForEach(numbers, func(number Token) error {
		resettableGearPtrs.Restore() // Restore to our previous save point

		for {
			part, err := resettableGearPtrs.Next()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}

			lineDiff := part.Line - number.Line
			if lineDiff > 1 {
				// No point in continuing to search for this part number as it's too far away
				return nil
			} else if lineDiff < -1 {
				// If the part is more than two lines above the number then we can skip it
				// but we can also save this position to restore for the next number as we
				// know it's too far away for that too
				resettableGearPtrs.Save()
			} else {
				numLeft := number.StartCol()
				numRight := number.EndCol()
				partLeft := part.StartCol() - 1
				partRight := part.EndCol() + 1

				if (numLeft <= partLeft && partLeft <= numRight) || (numLeft <= partRight && partRight <= numRight) {
					// log.Debug().Str("part_number", number.String()).Str("part_symbol", part.String()).Int("part_left", partLeft).Int("part_right", partRight).Int("num_left", numLeft).Int("num_right", numRight).Msg("Found overlap")
					part.Numbers = append(part.Numbers, number.Value)
					return nil
				}
			}
		}
	})
	if err != nil {
		return 0, err
	}

	resettableGearPtrs.Reset()
	gears := stream.Filter[*Token](resettableGearPtrs, func(part *Token) (bool, error) {
		return len(part.Numbers) == 2, nil
	})

	gearRatios := stream.Map(gears, func(part *Token) (int, error) {
		product := 1
		for _, number := range part.Numbers {
			num, err := strconv.Atoi(number)
			if err != nil {
				return 0, err
			}
			product *= num
		}
		return product, nil
	})

	return stream.Sum(gearRatios)
}
