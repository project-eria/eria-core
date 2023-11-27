package automations

import (
	"errors"
	"regexp"
	"strings"

	"github.com/project-eria/go-wot/consumer"
)

/*
 *
 * `context/<context name>`
 *	- not: `context/!<context name>`
 */
func contextCondition(conditionArray []string, contextsThing consumer.ConsumedThing) (bool, error) {
	// Check if the condition has the correct number of parameters
	if len(conditionArray) != 2 {
		return false, errors.New("invalid condition length")
	}

	// Check if the context name is valid
	match, _ := regexp.MatchString(`^!?\w*$`, conditionArray[1])
	if !match {
		return false, errors.New("invalid context name")
	}

	if strings.HasPrefix(conditionArray[1], "!") {
		// Check if the context is NOT active
		active, err := contextActive(contextsThing, conditionArray[1][1:])
		if err != nil {
			return false, err
		}
		return !active, nil
	} else {
		// Check if the context is active
		active, err := contextActive(contextsThing, conditionArray[1])
		if err != nil {
			return false, err
		}
		return active, nil
	}
}

func contextActive(contextsThing consumer.ConsumedThing, contextName string) (bool, error) {
	raw, err := contextsThing.ReadProperty(contextName)
	if err != nil {
		return false, err
	}
	if value, ok := raw.(bool); ok {
		return value, nil
	}
	return false, errors.New("invalid context value")
}
