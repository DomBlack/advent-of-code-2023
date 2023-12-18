package main

import (
	"time"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/pflag"

	// Register all the days
	_ "github.com/DomBlack/advent-of-code-2023/internal/day01"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day02"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day03"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day04"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day05"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day06"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day07"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day08"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day09"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day10"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day11"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day12"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day13"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day14"
)

func main() {
	var onlyDay int
	var verboseLevel int
	pflag.IntVarP(&onlyDay, "day", "d", 0, "Only run this day")
	pflag.CountVarP(&verboseLevel, "verbose", "v", "Increase verbosity")
	pflag.Parse()

	newLevel := zerolog.Level(int(zerolog.InfoLevel) - verboseLevel)
	if newLevel < zerolog.TraceLevel {
		newLevel = zerolog.TraceLevel
	}
	log.Logger = log.Level(newLevel)

	if onlyDay != 0 {
		log.Info().Int("day", onlyDay).Msg("Only running single day")
	} else {
		log.Info().Msg("Running all days")
	}

	days := runner.AllDays()
	runCount := 0

	start := time.Now()
	for _, day := range days {
		if onlyDay == 0 || day.Day() == onlyDay {
			day.Run()
			runCount++
		}
	}
	dur := time.Since(start)

	log.Info().Int("days", runCount).Str("duration", dur.String()).Msg("Finished running days")
}
