package automations

import (
	"errors"
	"strings"
	"time"

	"github.com/project-eria/go-wot/consumer"
)

// var schedules = map[time.Time][]Action{}

/**
 * Automation based on Time
 * The action is planned using gocron, and re-builded every day at 3:00 by default
 * (because of the summer/winter times changes), but that can be changed
 */

// TODO - re-build at specific time
// TODO - thing hour

/**
 * Handle time scheduling (hour)
 * `hour|hour|<time>`
 *	- fixed hour:
 *		- `at|hour|<hour in 15:04 format>`
 *		- `at|hour|<hour in 15:04:05 format>`
 *	- thing hour property: `at|hour|<thing device>:<property>`
 * 		- min option: `at|hour|<thing device>:<property>|min=<hour in 15:04:05 format>`
 *		- max option: `at|hour|<thing device>:<property>|max=<hour in 15:04:05 format>`
 * @param scheduleArray
 */
func atHourSchedule(scheduleArray []string, consumedThings map[string]consumer.ConsumedThing) (string, error) {
	if len(scheduleArray) < 3 || len(scheduleArray) > 5 {
		return "", errors.New("invalid at scheduling length")
	}

	// Test fixed hour: `at|hour|<hour in 15:04 format>`
	_, err := time.Parse("15:04", scheduleArray[2])
	if err == nil {
		// Correct time format
		return scheduleArray[2] + ":00", nil
	}
	// Test fixed hour: `at|hour|<hour in 15:04:05 format>`
	_, err = time.Parse("15:04:05", scheduleArray[2])
	if err == nil {
		// Correct time format
		return scheduleArray[2], nil
	}
	// Test thing hour property: `at|hour|<thing device>:<property>`
	timeThingArray := strings.Split(scheduleArray[2], ":")
	if len(timeThingArray) != 2 {
		return "", errors.New("invalid time")
	}
	// Does the time thing exist?
	timeThing, exists := consumedThings[timeThingArray[0]]
	if !exists {
		return "", errors.New("time thing doesn't exist")
	}

	// Get the property value
	timeThingValue, err := timeThing.ReadProperty(timeThingArray[1])
	if err != nil {
		return "", errors.New("time thing property not available: " + err.Error())
	}

	// Convert the thing date to a time, without the date/location, for comparison
	value, err := DateToHour(timeThingValue.(string))
	// TODO the format can be customised?
	if err != nil {
		return "", errors.New("invalid time thing value: " + err.Error())
	}

	// Test if min/max is set
	if len(scheduleArray) > 3 {
		for i := 3; i < len(scheduleArray); i++ {
			if strings.HasPrefix(scheduleArray[i], "min=") {
				min, err := time.Parse("15:04:05", scheduleArray[i][4:])
				if err != nil {
					return "", errors.New("invalid min time: " + err.Error())
				}
				if value.Before(min) {
					value = min
				}
			} else if strings.HasPrefix(scheduleArray[i], "max=") {
				max, err := time.Parse("15:04:05", scheduleArray[i][4:])
				if err != nil {
					return "", errors.New("invalid max time: " + err.Error())
				}
				if value.After(max) {
					value = max
				}
			}
		}
	}

	return value.Format("15:04:05"), nil
}

// func startTimeSchedules() {

// }
