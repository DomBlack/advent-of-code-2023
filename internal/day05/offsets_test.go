package day05

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapsToMergedOffset(t *testing.T) {
	input := `
seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4
`
	m, err := parseMaps([]byte(input))
	assert.NoError(t, err)

	assertMapping := func(seed int, what string, answer int, maps ...Map) {
		offsets, err := MapsToMergedOffset(maps...)
		assert.NoError(t, err)

		idx := sort.Search(len(offsets), func(i int) bool {
			return offsets[i].InputTo >= seed
		})

		lastAnswer := seed
		if idx < len(offsets) {
			lastAnswer += offsets[idx].OffsetBy
		}

		assert.Equal(t, answer, lastAnswer, fmt.Sprintf("expected %s %d, got %s %d", what, answer, what, lastAnswer))
	}

	assertMapping(82, "soil", 84, m.SeedToSoil)
	assertMapping(82, "fertilizer", 84, m.SeedToSoil, m.SoilToFertilizer)
	assertMapping(82, "water", 84, m.SeedToSoil, m.SoilToFertilizer, m.FertilizerToWater)
	assertMapping(82, "light", 77, m.SeedToSoil, m.SoilToFertilizer, m.FertilizerToWater, m.WaterToLight)
	assertMapping(82, "temperature", 45, m.SeedToSoil, m.SoilToFertilizer, m.FertilizerToWater, m.WaterToLight, m.LightToTemperature)
	assertMapping(82, "humidity", 46, m.SeedToSoil, m.SoilToFertilizer, m.FertilizerToWater, m.WaterToLight, m.LightToTemperature, m.TemperatureToHumidity)
	assertMapping(82, "location", 46, m.SeedToSoil, m.SoilToFertilizer, m.FertilizerToWater, m.WaterToLight, m.LightToTemperature, m.TemperatureToHumidity, m.HumidityToLocation)
}
