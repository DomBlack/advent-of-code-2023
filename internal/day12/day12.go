package day12

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day12 = runner.NewStreamingDay(12, parseInput, part1, part2).
	WithExpectedAnswers("7857", "28606137449920")

func part1(log zerolog.Logger, input stream.Stream[Springs]) (answer string, err error) {
	cache := make(map[string]int)

	combinations := stream.Map(input, func(spring Springs) (int, error) {
		return countOptions(cache, spring.Springs, spring.DamagedSpringGroups)
	})

	return stream.SumToString(combinations)
}

func part2(log zerolog.Logger, input stream.Stream[Springs]) (answer string, err error) {
	cache := make(map[string]int)

	combinations := stream.Map(input, func(spring Springs) (int, error) {
		spring = spring.Unfold()
		return countOptions(cache, spring.Springs, spring.DamagedSpringGroups)
	})

	return stream.SumToString(combinations)
}

func countOptions(cache map[string]int, spring SpringConditionList, groups []int) (numOptions int, err error) {
	// Cache the results
	key := fmt.Sprintf("%s-%v", spring.String(), groups)
	if val, ok := cache[key]; ok {
		return val, nil
	}
	defer func() {
		if err == nil {
			cache[key] = numOptions
		}
	}()

	// If we've reached the end of the spring list then we're done
	if len(spring) == 0 {
		if len(groups) == 0 {
			return 1, err
		} else {
			return 0, err
		}
	}

	// Short circuit if we have more groups than springs
	groupSum := 0
	for _, group := range groups {
		groupSum += group
	}
	if groupSum > len(spring) {
		return 0, err
	}

	// Now check the springs
	switch spring[0] {
	case Unknown:
		// recurse with the next spring set to both states and the same groups
		operation, err := countOptions(cache, append(SpringConditionList{Operational}, spring[1:]...), groups)
		if err != nil {
			return 0, err
		}

		damaged, err := countOptions(cache, append(SpringConditionList{Damaged}, spring[1:]...), groups)
		if err != nil {
			return 0, err
		}
		return operation + damaged, nil

	case Operational:
		// recurse with the next spring and the same groups
		return countOptions(cache, spring[1:], groups)

	case Damaged:
		// Here we do the actual hard work

		if len(groups) == 0 {
			// There should be no more damaged springs
			return 0, err
		}

		// Ensure there isn't a working spring before this damaged spring group finishes
		for i := 0; i < groups[0]; i++ {
			if spring[i] == Operational {
				return 0, err
			}
		}

		// We can now assume the groups[0] is all damanged

		if len(spring) > groups[0] {
			if spring[groups[0]] == Damaged {
				// If the next spring is also damaged then this group is too big
				return 0, err
			}

			return countOptions(cache, spring[groups[0]+1:], groups[1:])
		} else {
			return countOptions(cache, nil, groups[1:])
		}

	default:
		return 0, errors.Newf("unknown spring condition %q", spring[0])
	}
}

type SpringCondition uint8

const (
	Unknown SpringCondition = iota
	Operational
	Damaged
)

type SpringConditionList []SpringCondition

func (s SpringConditionList) String() string {
	var sb strings.Builder

	for _, spring := range s {
		switch spring {
		case Unknown:
			sb.WriteRune('?')
		case Operational:
			sb.WriteRune('.')
		case Damaged:
			sb.WriteRune('#')
		}
	}

	return sb.String()
}

func (s SpringConditionList) ValidAgainst(damagedGroups []int, partial bool) bool {
	count := 0
	groupIdx := 0

	for _, spring := range s {
		if spring == Damaged {
			if groupIdx >= len(damagedGroups) {
				return false
			}

			count++
		} else if count > 0 {
			if count != damagedGroups[groupIdx] {
				return false
			}
			count = 0
			groupIdx++
		}
	}

	if partial {
		// If we have a partial list then we need to check that the last group is valid
		// i.e. it's not already bigger than the damaged group it would be matching
		if count > 0 && len(damagedGroups) > groupIdx && damagedGroups[groupIdx] < count {
			return false
		}
	} else if count > 0 {
		// This is meant to be the last group and it has the correct number of damaged springs
		return damagedGroups[groupIdx] == count && groupIdx == len(damagedGroups)-1
	} else if groupIdx != len(damagedGroups) {
		// not enough groups
		return false
	}

	return true
}

type Springs struct {
	Springs             SpringConditionList
	DamagedSpringGroups []int
}

func (s Springs) Unfold() Springs {
	const factor = 5

	length := len(s.Springs) + 1

	rtn := Springs{
		Springs:             make(SpringConditionList, length*factor-1),
		DamagedSpringGroups: make([]int, len(s.DamagedSpringGroups)*factor),
	}

	for i := 0; i < factor; i++ {
		copy(rtn.Springs[i*length:], s.Springs)
		copy(rtn.DamagedSpringGroups[i*len(s.DamagedSpringGroups):], s.DamagedSpringGroups)
	}

	return rtn
}

func (s Springs) String() string {
	var sb strings.Builder

	sb.WriteString(s.Springs.String())

	sb.WriteRune(' ')

	for i, group := range s.DamagedSpringGroups {
		if i != 0 {
			sb.WriteRune(',')
		}

		sb.WriteString(strconv.Itoa(group))
	}

	return sb.String()
}

func parseInput(input []byte) stream.Stream[Springs] {
	return stream.Map(stream.LinesFrom(input), func(line string) (Springs, error) {
		springs, groups, found := strings.Cut(line, " ")
		if !found {
			return Springs{}, errors.Newf("could not find seperator between springs and groups for %q", line)
		}

		spring := Springs{}

		// Parse the spring conditions
		for _, c := range springs {
			switch c {
			case '?':
				spring.Springs = append(spring.Springs, Unknown)
			case '.':
				spring.Springs = append(spring.Springs, Operational)
			case '#':
				spring.Springs = append(spring.Springs, Damaged)
			default:
				return Springs{}, errors.Newf("unknown spring condition %q", c)
			}
		}

		// Parse the spring groups
		for _, g := range strings.Split(groups, ",") {
			num, err := strconv.Atoi(g)
			if err != nil {
				return Springs{}, errors.Wrapf(err, "could not parse spring group %q", g)
			}

			spring.DamagedSpringGroups = append(spring.DamagedSpringGroups, num)
		}

		return spring, nil
	})
}
