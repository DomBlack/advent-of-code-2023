package day05

import (
	"math"
	"slices"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/rs/zerolog"
)

var Day05 = runner.NewDay(5, parseMaps, part1, part2).
	WithExpectedAnswers(226172555, 47909639)

func part1(log zerolog.Logger, input Maps) (answer int, err error) {
	lowest := math.MaxInt

	for _, seed := range input.Seeds {
		soil := input.SeedToSoil.Destination(seed)
		fertilizer := input.SoilToFertilizer.Destination(soil)
		water := input.FertilizerToWater.Destination(fertilizer)
		light := input.WaterToLight.Destination(water)
		temp := input.LightToTemperature.Destination(light)
		humidity := input.TemperatureToHumidity.Destination(temp)
		location := input.HumidityToLocation.Destination(humidity)

		if lowest > location {
			lowest = location
		}
	}

	return lowest, nil
}

func part2(log zerolog.Logger, input Maps) (answer int, err error) {
	offsets, err := MapsToMergedOffset(
		input.SeedToSoil,
		input.SoilToFertilizer,
		input.FertilizerToWater,
		input.WaterToLight,
		input.LightToTemperature,
		input.TemperatureToHumidity,
		input.HumidityToLocation,
	)
	if err != nil {
		return 0, err
	}

	seedToLocations := offsets.ToMap()

	seedRanges := make([]Range, 0)
	for i := 0; i < len(input.Seeds); i += 2 {
		seedRanges = append(seedRanges, Range{
			SrcRangeStart:  input.Seeds[i],
			DestRangeStart: input.Seeds[i],
			Length:         input.Seeds[i+1],
		})
	}
	slices.SortFunc(seedRanges, func(a, b Range) int {
		return a.SrcRangeStart - b.SrcRangeStart
	})

	lowest := math.MaxInt
	for _, toLocation := range seedToLocations {
		for _, seedRange := range seedRanges {

			if toLocation.SourcesOverLap(seedRange) {
				if toLocation.DestRangeStart < lowest {
					lowest = toLocation.DestRangeStart
				}
				break
			}
		}
	}

	return lowest, nil
}
