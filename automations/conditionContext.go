package automations

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type conditionContext struct {
	context string
	invert  bool
}

/*
 * `context/<context name>`
 *	- not: `context/!<context name>`
 */
func NewConditionContext(conditionArray []string) (*conditionContext, error) {
	if _contextsThing == nil {
		return nil, errors.New("contexts thing not configured")
	}
	// Check if the condition has the correct number of parameters
	if len(conditionArray) != 2 {
		return nil, errors.New("invalid condition length")
	}

	// Check if the context name is valid
	match, _ := regexp.MatchString(`^!?\w*$`, conditionArray[1])
	if !match {
		return nil, errors.New("invalid context name")
	}

	if strings.HasPrefix(conditionArray[1], "!") {
		// If the context is NOT active
		return &conditionContext{
			context: conditionArray[1][1:],
			invert:  true,
		}, nil
	} else {
		// If the context is active
		return &conditionContext{
			context: conditionArray[1],
			invert:  false,
		}, nil
	}
}

func (condition *conditionContext) check(time.Time) (bool, error) {
	active, err := contextActive(condition.context)
	if err != nil {
		return false, err
	}
	return condition.invert != active, nil
}

func contextActive(contextName string) (bool, error) {
	raw, err := _contextsThing.ReadProperty(contextName)
	if err != nil {
		return false, err
	}
	if value, ok := raw.(bool); ok {
		return value, nil
	}
	return false, errors.New("invalid context value")
}
