package automations

import (
	"errors"
	"strings"
	"time"

	"github.com/gookit/goutil/arrutil"
	zlog "github.com/rs/zerolog/log"
)

var (
	newContextCondition  = NewConditionContext
	newTimeCondition     = NewConditionTime
	newPropertyCondition = NewConditionProperty
) // Mocking for inner functions

type Condition interface {
	check(time.Time) (bool, error)
}

type Observables struct {
	contexts   []string
	properties map[string][]string // properties to observe, by thing ref
}

func getConditions(conditions []string, thingRefs []string) ([]Condition, *Observables, error) {
	observables := &Observables{
		contexts:   []string{},
		properties: map[string][]string{},
	}
	cs := make([]Condition, 0)
	if conditions == nil {
		zlog.Trace().Msg("[automations:getConditions] no conditions")
		return cs, nil, nil
	}
	for _, condition := range conditions {
		conditionArray := strings.Split(condition, "|")
		var err error
		var c Condition
		switch conditionArray[0] {
		case "context":
			c, err = newContextCondition(conditionArray)
			// Add the context to the observables list
			if c.(*conditionContexts) != nil {
				// Add the context to the observables list, and remove duplicates
				observables.contexts = arrutil.Union(observables.contexts, c.(*conditionContexts).list(), arrutil.StringEqualsComparer)

			}
		case "time":
			c, err = newTimeCondition(conditionArray)
		case "property":
			var pc *conditionProperty
			pc, err = newPropertyCondition(conditionArray, thingRefs)
			if pc != nil {
				// Add the properties to the observables list, and remove duplicates
				for thingRef, properties := range pc.list() {
					observables.properties[thingRef] = arrutil.Union(observables.properties[thingRef], properties, arrutil.StringEqualsComparer)
				}
			}
			c = pc
		default:
			return nil, nil, errors.New("invalid condition type")
		}
		if err != nil {
			return nil, nil, err
		}
		cs = append(cs, c)
	}
	return cs, observables, nil
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
