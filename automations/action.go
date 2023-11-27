package automations

import (
	"errors"
	"strings"

	"github.com/project-eria/go-wot/producer"
)

type Action struct {
	Ref        string
	Parameters map[string]interface{}
	Handler    producer.ActionHandler
}

/**
 * `<action>/<param name>=<value>/<param name>=<value>`
 */
func getAction(exposedThing producer.ExposedThing, actionArray []string) (Action, error) {
	// Check if the action has the correct number of parameters
	if len(actionArray) == 0 {
		return Action{}, errors.New("invalid action length")
	}
	ref := actionArray[0]

	// Check if the action exists
	action, err := exposedThing.ExposedAction(ref)

	if err != nil {
		return Action{}, err
	}

	// Get the parameters
	parameters := make(map[string]interface{})
	for i := 1; i < len(actionArray); i++ {
		param := strings.Split(actionArray[i], "=")
		paramName := param[0]
		paramValue := param[1]
		parameters[paramName] = paramValue
	}

	return Action{
		Ref:        ref,
		Handler:    action.GetHandler(),
		Parameters: parameters,
	}, nil
}
