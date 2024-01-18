package automations

import (
	"errors"
)

type scheduleNone struct {
}

/**
 * `none`: Disable the scheduling
 */
func NewScheduleNone(scheduleArray []string) (*scheduleNone, error) {
	// Check if the schedule has the correct number of parameters
	if len(scheduleArray) != 1 {
		return nil, errors.New("invalid none schedule length")
	}
	return &scheduleNone{}, nil
}

func (s *scheduleNone) start(action ActionRunner) error {
	if action == nil {
		return errors.New("missing action")
	}
	return nil
}

func (s *scheduleNone) job() error {
	// nothing to do
	return nil
}

func (s *scheduleNone) cancel() {
	// nothing to do
}

func (s *scheduleNone) equals(other Schedule) bool {
	if other == nil {
		return false
	}
	_, ok := other.(*scheduleNone)

	return ok && true // is always true as the only way to have a different immediate schedule is from nothing
}

func (s *scheduleNone) string() string {
	return "disabled"
}
