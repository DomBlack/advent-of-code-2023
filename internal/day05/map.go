package day05

import (
	"fmt"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
)

type Maps struct {
	Seeds                 []int
	SeedToSoil            Map
	SoilToFertilizer      Map
	FertilizerToWater     Map
	WaterToLight          Map
	LightToTemperature    Map
	TemperatureToHumidity Map
	HumidityToLocation    Map
}

type Map []Range

func (m Map) String() string {
	var str strings.Builder

	for _, r := range m {
		str.WriteRune('\t')
		str.WriteString(r.String())
		str.WriteRune('\n')
	}

	return str.String()
}

type Range struct {
	DestRangeStart int
	SrcRangeStart  int
	Length         int
}

// SourcesOverLap returns true if these two ranges contain any overlap at all?
//
// i.e. does any part of `r`'s SrcRangeStart -> (SrcRangeStart + Length)
// overlap with `other`'s SrcRangeStart -> (SrcRangeStart + Length)
func (r Range) SourcesOverLap(other Range) bool {
	return r.SrcRangeStart < other.SrcRangeStart+other.Length &&
		r.SrcRangeStart+r.Length > other.SrcRangeStart
}

func (r Range) String() string {
	return fmt.Sprintf("Range(%d->%d => %d->%d)", r.SrcRangeStart, r.SrcRangeStart+r.Length-1, r.DestRangeStart, r.DestRangeStart+r.Length-1)
}

func (m Map) Destination(source int) int {
	// Find the smallest range which starts after the source
	i := sort.Search(len(m), func(i int) bool {
		return m[i].SrcRangeStart > source
	})

	if i > len(m) || i <= 0 {
		// No matches; source maps to destination
		return source
	}

	rng := m[i-1]
	offset := source - rng.SrcRangeStart
	if offset < 0 || offset >= rng.Length {
		// Outside the range
		return source
	}

	return rng.DestRangeStart + offset
}

func parseMaps(input []byte) (rtn Maps, err error) {
	lines := stream.LinesFrom(input)

	// Read the seeds line
	seeds, err := lines.Next()
	if err != nil {
		return Maps{}, errors.Wrap(err, "unable to read seeds line")
	} else if !strings.HasPrefix(seeds, "seeds: ") {
		return Maps{}, errors.Newf("seeds line has unexpected prefix: %q", seeds)
	}
	for _, str := range strings.Fields(seeds[7:]) {
		seed, err := strconv.Atoi(str)
		if err != nil {
			return Maps{}, errors.Wrapf(err, "unable to parse seed: %q", str)
		}

		rtn.Seeds = append(rtn.Seeds, seed)
	}

	// Read the expected blank line after seeds
	blank, err := lines.Next()
	if err != nil {
		return Maps{}, errors.Wrap(err, "unable to read blank line")
	} else if blank != "" {
		return Maps{}, errors.Newf("unexpected text on line: %q", blank)
	}

	for {
		// Parse the map name
		mapName, err := lines.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return Maps{}, errors.Wrap(err, "unexpected error reading map name")
		} else if !strings.HasSuffix(mapName, " map:") {
			return Maps{}, errors.Newf("line was missing map suffix: %q", mapName)
		}
		mapName = strings.TrimSuffix(mapName, " map:")

		mapValue, err := parseMap(lines)
		if err != nil {
			return Maps{}, errors.Wrapf(err, "error parsing map %q", mapName)
		}

		switch mapName {
		case "seed-to-soil":
			rtn.SeedToSoil = mapValue
		case "soil-to-fertilizer":
			rtn.SoilToFertilizer = mapValue
		case "fertilizer-to-water":
			rtn.FertilizerToWater = mapValue
		case "water-to-light":
			rtn.WaterToLight = mapValue
		case "light-to-temperature":
			rtn.LightToTemperature = mapValue
		case "temperature-to-humidity":
			rtn.TemperatureToHumidity = mapValue
		case "humidity-to-location":
			rtn.HumidityToLocation = mapValue
		default:
			return Maps{}, errors.Newf("unexpected map name %q", mapName)
		}
	}

	return rtn, nil
}

func parseMap(lines stream.Stream[string]) (rtn Map, err error) {
	// Parse each range
	for {
		line, err := lines.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "unable to parse map")
		}

		if strings.TrimSpace(line) == "" {
			break
		}

		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, errors.Newf("line had unexpected number of fields (%d) %q", len(parts), line)
		}

		r := Range{}
		r.DestRangeStart, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse dest range start from %q", parts[0])
		}

		r.SrcRangeStart, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse src range start from %q", parts[1])
		}

		r.Length, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse length from %q", parts[2])
		}

		rtn = append(rtn, r)
	}

	slices.SortFunc(rtn, func(a, b Range) int {
		return a.SrcRangeStart - b.SrcRangeStart
	})

	return rtn, nil
}
