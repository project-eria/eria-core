package automations

import "errors"

/**
 * `<empty>`: Immediate run
 */
func immediateSchedule(scheduleArray []string) (string, error) {
	// Check if the schedule has the correct number of parameters
	if len(scheduleArray) != 1 {
		return "", errors.New("invalid immediate schedule length")
	}
	return "", nil
}
