package automations

import (
	"errors"
	"strings"
	"time"

	zlog "github.com/rs/zerolog/log"
)

var (
	newContextCondition = NewConditionContext
	newTimeCondition    = NewConditionTime
) // Mocking for inner functions

type Condition interface {
	check(time.Time) (bool, error)
}

func getConditions(conditions []string) ([]Condition, error) {
	cs := make([]Condition, 0)
	if conditions == nil {
		zlog.Trace().Msg("[automations:getConditions] no conditions")
		return cs, nil
	}
	for _, condition := range conditions {
		conditionArray := strings.Split(condition, "|")
		var err error
		var c Condition
		switch conditionArray[0] {
		case "context":
			c, err = newContextCondition(conditionArray)
		case "time":
			c, err = newTimeCondition(conditionArray)
		// TODO case "property":
		// 	c, err = NewConditionProperty(conditionArray)
		default:
			return nil, errors.New("invalid condition type")
		}
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func checkConditions(conditions []Condition, now time.Time) (bool, error) {
	for _, condition := range conditions {
		ok, err := condition.check(now)
		if err != nil {
			return false, err
		} else if !ok {
			// If one condition is not true, we skip the group
			return false, nil
		}
	}
	return true, nil
}
