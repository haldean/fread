package main

import (
    "fmt"
	"time"
)

func DateFormat(t time.Time) string {
	diff := time.Now().Sub(t)

	var val uint
	var unit string

	if diff >= 24 * time.Hour {
		val = uint(diff.Hours() / 24)
		if val > 1 {
			unit = "days"
		} else {
			unit = "day"
		}
	} else if diff >= time.Hour {
		val = uint(diff.Hours())
		if val > 1 {
			unit = "hours"
		} else {
			unit = "hour"
		}
	} else if diff >= time.Minute {
		val = uint(diff.Minutes())
		if val > 1 {
			unit = "minutes"
		} else {
			unit = "minute"
		}
	} else if diff >= time.Second {
		val = uint(diff.Seconds())
		if val > 1 {
			unit = "seconds"
		} else {
			unit = "second"
		}
	} else {
		val = uint(diff.Seconds() / 1000)
		unit = "ms"
	}
	return fmt.Sprintf("%3d %-7s", val, unit)
}

func ShowHeadlines(t Timeline) {
    for _, item := range t {
        fmt.Printf("%v   %v\n", DateFormat(item.Date), item.Title)
    }
}
