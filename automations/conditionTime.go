package automations

import (
	"errors"
	"strings"
	"time"
)

type conditionTime struct {
	before string
	after  string
}

/**
 * - `time|[before/after]=<hour in 15:04/15:04:05 format>`
 */
func NewConditionTime(conditionArray []string) (*conditionTime, error) {
	if len(conditionArray) == 2 { // Check if the condition has the correct number of parameters
		ba := strings.Split(conditionArray[1], "=")
		if len(ba) != 2 {
			return nil, errors.New("invalid condition parameter")
		}
		_, err := time.Parse("15:04", ba[1])
		if err != nil {
			return nil, errors.New("invalid condition time")
		}
		switch ba[0] {
		case "after":
			return &conditionTime{after: ba[1]}, nil
		case "before":
			return &conditionTime{before: ba[1]}, nil
		}
		return nil, errors.New("invalid condition parameter type")
	}
	return nil, errors.New("invalid condition length")
}

// HourToCurrentTime needs to be used on checking
func (c *conditionTime) check(now time.Time) (bool, error) {
	if c.after != "" {
		t, err := HourToCurrentTime(c.after, now) // Convert the hour string to a time
		if err != nil {
			return false, err
		}
		return now.After(t), nil
	} else if c.before != "" {
		t, err := HourToCurrentTime(c.before, now) // Convert the hour string to a time
		if err != nil {
			return false, err
		}
		return now.Before(t), nil
	}
	return false, errors.New("unexpected invalid condition time")
}
