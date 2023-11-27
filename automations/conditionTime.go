package automations

import (
	"errors"
	"strings"
	"time"
)

/**
 * - `time|[before/after]=<hour in 15:00 format>`
 */
func timeCondition(conditionArray []string, now time.Time) (bool, error) {
	if len(conditionArray) == 2 { // Check if the condition has the correct number of parameters
		ba := strings.Split(conditionArray[1], "=")
		if len(ba) != 2 {
			return false, errors.New("invalid condition parameter")
		}
		t, err := HourToTime(ba[1], now) // Convert the hour string to a time
		if err == nil {
			switch ba[0] {
			case "after":
				return now.After(t), nil
			case "before":
				return now.Before(t), nil
			}
			return false, errors.New("invalid condition parameter type")
		}
		return false, errors.New("invalid condition time")
	}
	return false, errors.New("invalid condition length")
}
