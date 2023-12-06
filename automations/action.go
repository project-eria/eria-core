package automations

import (
	"errors"
	"strings"

	"github.com/project-eria/go-wot/producer"
)

type Action interface {
	run() error
}

type action struct {
	Ref        string
	Parameters map[string]interface{}
	Handler    producer.ActionHandler
}

/**
 * `<action>/<param name>=<value>/<param name>=<value>`
 */
func getAction(exposedThing producer.ExposedThing, actionStr string) (*action, error) {
	actionStr = strings.TrimSpace(actionStr)
	if actionStr == "" {
		return nil, errors.New("missing action configuration") // Skip this automation
	}
	actionArray := strings.Split(actionStr, "|")
	ref := actionArray[0]

	// Check if the action exists
	exposedAction, err := exposedThing.ExposedAction(ref)
	if err != nil {
		return nil, err
	}

	// Get the parameters
	parameters := make(map[string]interface{})
	for i := 1; i < len(actionArray); i++ {
		param := strings.Split(actionArray[i], "=")
		paramName := param[0]
		paramValue := param[1]
		parameters[paramName] = paramValue
	}

	return &action{
		Ref:        ref,
		Handler:    exposedAction.GetHandler(),
		Parameters: parameters,
	}, nil
}

func (a *action) run() error {
	// TODO
	return nil
}
