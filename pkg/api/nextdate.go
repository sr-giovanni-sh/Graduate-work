package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	FormatDate    = "20060102"
	maxIterations = 2000
)

/*
NextDate calculates the next occurrence date based on a start date and a repetition rule.

Parameters:
- now: the reference time used to determine the next future date.
- dstart: the start date string in FormatDate format.
- repeat: the repetition rule string that defines how often the event repeats.

Repetition rule format:
- "d <interval>": daily with a given interval in days (1–400).
- "y": yearly (adds 1 year at a time).
- "w <days>": weekly with specific weekdays (1–7, comma‑separated).
- "m <days>[,<months>]": monthly with specific days of the month and optional specific months.
  - Days can be positive (1–31), -1 for the last day of the month, or -2 for the second‑to‑last day.
  - Months are 1–12; if omitted, all months are allowed.

Behavior:
  - Validates that the repeat rule is not empty and that the start date is parseable.
  - For each rule type, iteratively advances the date until it finds a date strictly after `now`.
  - Enforces a maximum number of iterations (maxIterations) to prevent infinite loops.
  - Returns the resulting date formatted as a string in FormatDate, or an error if validation
    fails, the rule is unknown, or the iteration limit is reached.
*/
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat rule cannot be empty")
	}

	date, err := time.Parse(FormatDate, dstart)
	if err != nil {
		return dstart, fmt.Errorf("error parse date: %v", err)
	}

	parts := strings.Fields(repeat)
	firstLatter := parts[0]

	switch firstLatter {
	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("missing interval for rule 'd'")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", fmt.Errorf("invalid interval for rule 'd'")
		}

		if now.After(date) {
			daysDiff := int(now.Sub(date).Hours() / 24)
			if daysDiff > interval {
				skipIntervals := daysDiff / interval
				date = date.AddDate(0, 0, skipIntervals*interval)
			}
		}

		iterations := 0
		for {
			iterations++
			if iterations > maxIterations {
				return "", fmt.Errorf("max iterations reached for rule 'd'")
			}

			date = date.AddDate(0, 0, interval)
			if afterDay(date, now) {
				break
			}
		}

		return date.Format(FormatDate), nil

	case "y":
		iterations := 0
		for {
			iterations++
			if iterations > maxIterations {
				return "", fmt.Errorf("max iterations reached for rule 'y'")
			}

			date = date.AddDate(1, 0, 0)
			if afterDay(date, now) {
				break
			}
		}

		return date.Format(FormatDate), nil

	case "w":
		if len(parts) < 2 {
			return "", fmt.Errorf("missing days for rule 'w'")
		}

		dateWorks := strings.Split(parts[1], ",")
		targetDays := make(map[int]bool)

		for _, dayStr := range dateWorks {
			dayNum, err := strconv.Atoi(dayStr)
			if err != nil || dayNum < 1 || dayNum > 7 {
				return "", fmt.Errorf("invalid weekday: %s", dayStr)
			}
			targetDays[dayNum] = true
		}

		iterations := 0
		for {
			iterations++
			if iterations > maxIterations {
				return "", fmt.Errorf("max iterations reached for rule 'w'")
			}

			date = date.AddDate(0, 0, 1)

			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}

			if targetDays[weekday] && date.After(now) {
				break
			}
		}

		return date.Format(FormatDate), nil

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("missing days for rule 'm'")
		}

		targetDays := make(map[int]bool)
		targetMonths := make(map[int]bool)

		dayStrings := strings.Split(parts[1], ",")
		for _, dayStr := range dayStrings {
			day, err := strconv.Atoi(dayStr)
			if err != nil || (day < 1 && day != -1 && day != -2) || day > 31 {
				return "", fmt.Errorf("invalid day value in rule 'm': %s", dayStr)
			}
			targetDays[day] = true
		}

		if len(parts) == 3 {
			monthStrings := strings.Split(parts[2], ",")
			for _, monthStr := range monthStrings {
				month, err := strconv.Atoi(monthStr)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("invalid month value in rule 'm': %s", monthStr)
				}
				targetMonths[month] = true
			}

		} else {
			for m := 1; m <= 12; m++ {
				targetMonths[m] = true
			}
		}

		iterations := 0
		for {
			iterations++
			if iterations > maxIterations {
				return "", fmt.Errorf("max iterations reached for rule 'm'")
			}

			date = date.AddDate(0, 0, 1)

			currentMonth := int(date.Month())
			if !targetMonths[currentMonth] {
				continue
			}

			firstOfNextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
			lastDayOfCurrentMonth := firstOfNextMonth.AddDate(0, 0, -1).Day()

			currentDay := date.Day()
			dayMatches := false

			if targetDays[currentDay] ||
				(targetDays[-1] && currentDay == lastDayOfCurrentMonth) ||
				(targetDays[-2] && currentDay == lastDayOfCurrentMonth-1) {
				dayMatches = true
			}

			if dayMatches && date.After(now) {
				break
			}
		}
		return date.Format(FormatDate), nil

	default:
		return "", fmt.Errorf("unknown repeat rule: %s", firstLatter)
	}
}

/*
afterDay checks whether a given date is strictly after another date.

Returns true if `date` occurs after `now`, otherwise false.
*/
func afterDay(date, now time.Time) bool {
	return date.After(now)
}
