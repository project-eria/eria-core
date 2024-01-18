package automations

import (
	"errors"
	"strings"
)

type Schedule interface {
	job() error
	start(ActionRunner) error
	equals(Schedule) bool
	cancel()
	string() string
}

var (
	newImmediateSchedule = NewScheduleImmediate
	newAtHourSchedule    = NewScheduleAtHour
	newNoneSchedule      = NewScheduleNone
) // Mocking for inner functions

func getSchedule(schedule string) (Schedule, error) {
	scheduleArray := strings.Split(schedule, "|")
	switch scheduleArray[0] {
	case "immediate":
		return newImmediateSchedule(scheduleArray)
	case "at":
		if scheduleArray[1] == "hour" {
			return newAtHourSchedule(scheduleArray)
			// TODO } else if scheduleArray[1] == "date" {
		} else {
			return nil, errors.New("invalid 'at' schedule type")
		}
	case "none":
		return newNoneSchedule(scheduleArray)
		// TODO case "every":
		// TODO case: "in":
	}
	return nil, errors.New("invalid schedule type")
}
