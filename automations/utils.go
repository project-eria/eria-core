package automations

import (
	"time"
)

/**
 * Convert hour string to time, with the location of the current time
 */
func HourToCurrentTime(hourStr string, now time.Time) (time.Time, error) {
	t, err := HourToTime(hourStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location()), nil
}

/**
* Convert hour string to time, without the date/location, for comparison
* Example: "11:12:13" -> time.Date(0, time.January, 1, 11, 12, 13, 0, time.UTC)
 */
func HourToTime(hourStr string) (time.Time, error) {
	var t time.Time
	var err error
	t, err = time.Parse("15:04", hourStr)
	if err != nil {
		t, err = time.Parse("15:04:05", hourStr)
	}
	return t, err
}

/**
 * Convert date to hour, without the date/location, for comparison
 * Example: "2024-03-02T11:12:13" -> time.Date(0, time.January, 1, 11, 12, 13, 0, time.UTC)
 */
func DateToHour(dateStr string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	hourStr := date.Format("15:04:05")
	t, _ := time.Parse("15:04:05", hourStr)
	return t, nil
}
