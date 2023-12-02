package main

import (
	"time"

	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/rs/zerolog/log"

	"github.com/spf13/pflag"

	// Register all the days
	_ "github.com/DomBlack/advent-of-code-2023/internal/day01"
	_ "github.com/DomBlack/advent-of-code-2023/internal/day02"
)

func main() {
	var onlyDay int
	pflag.IntVar(&onlyDay, "day", 0, "Only run this day")

	pflag.Parse()

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
