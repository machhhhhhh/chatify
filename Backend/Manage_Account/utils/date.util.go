package utils

import (
	"strconv"
	"time"
)

func GetCurrentYear() string {
	// Convert the current year
	return strconv.Itoa(time.Now().Year())
}

func ConvertDate(input string) (time.Time, error) {
	const layout string = "2006-01-02 15:04:05"
	return time.Parse(layout, input)
}

func StartOfDay(input time.Time) time.Time {
	year, month, day := input.Date()
	// midnight at the start of that day
	return time.Date(year, month, day, 0, 0, 0, 0, input.Location())
}

func EndOfDay(input time.Time) time.Time {
	// start of next day minus 1 nanosecond
	return StartOfDay(input).Add(24*time.Hour - time.Nanosecond)
}

func EndOfYesterday() time.Time {
	// Subtract 1 day from today and get the end of that day
	return EndOfDay(time.Now().AddDate(0, 0, -1))
}
