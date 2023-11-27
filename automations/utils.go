package automations

import "time"

/**
 * Convert hour string to time, with the location of the current time
 */
func HourToTime(hourStr string, now time.Time) (time.Time, error) {
	hour, err := time.Parse("15:04", hourStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(now.Year(), now.Month(), now.Day(), hour.Hour(), hour.Minute(), 0, 0, now.Location()), nil
}

/**
 * Convert date to hour, without the date/location, for comparison
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
