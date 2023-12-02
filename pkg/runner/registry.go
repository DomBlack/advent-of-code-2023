package runner

import (
	"slices"
)

var (
	days = make(map[int]RunnableDay)
)

type RunnableDay interface {
	// Day returns the day number
	Day() int

	// Run runs the day
	Run()
}

// AllDays returns all the days in order
func AllDays() []RunnableDay {
	var daysSlice []RunnableDay

	for _, day := range days {
		daysSlice = append(daysSlice, day)
	}

	slices.SortFunc(daysSlice, func(a, b RunnableDay) int {
		return a.Day() - b.Day()
	})

	return daysSlice
}
