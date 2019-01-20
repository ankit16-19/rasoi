package main

import (
	"time"
)

// GetDateFromTime :
func GetDateFromTime(t time.Time) string {
	return t.Format(time.RFC3339)[:10]
}

// FirstDayofWeek :
func FirstDayofWeek(t time.Time) time.Time {
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, -1)
	}
	return t
}

// WholeWeekDates :
func WholeWeekDates(t time.Time) []time.Time {
	var array []time.Time
	t = FirstDayofWeek(t)
	array = append(array, t)
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, 1)
		array = append(array, t)
	}
	return array
}
