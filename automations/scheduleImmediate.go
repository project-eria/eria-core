package automations

import (
	"errors"
)

type scheduleImmediate struct {
}

/**
 * `immediate`: Immediate run
 */
// TODO should not activate on start, only on condition switch
func NewScheduleImmediate(scheduleArray []string) (*scheduleImmediate, error) {
	// Check if the schedule has the correct number of parameters
	if len(scheduleArray) != 1 {
		return nil, errors.New("invalid immediate schedule length")
	}
	return &scheduleImmediate{}, nil
}

func (s *scheduleImmediate) start(action Action) error {
	if action == nil {
		return errors.New("missing action")
	}
	return action.run()
}

func (s *scheduleImmediate) job() error {
	// nothing to do
	return nil
}

func (s *scheduleImmediate) cancel() {
	// nothing to do
}

func (s *scheduleImmediate) equals(other Schedule) bool {
	if other == nil {
		return false
	}
	_, ok := other.(*scheduleImmediate)

	return ok && true // is always true as the only way to have a different immediate schedule is from nothing
}

func (s *scheduleImmediate) string() string {
	return "immediately"
}
