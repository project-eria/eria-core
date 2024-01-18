package automations

import (
	"errors"
	"strings"

	"github.com/project-eria/go-wot/producer"
)

type ActionRunner interface {
	run() error
}

type Action struct {
	AutomationName string
	Ref            string
	ExposedThing   producer.ExposedThing
	Value          interface{}
	Parameters     map[string]string
	ExposedAction  producer.ExposedAction
}

/**
 * `<action>|<value>|<param name>=<value>|<param name>=<value>`
 */
func getAction(exposedThing producer.ExposedThing, automationName string, actionStr string) (*Action, error) {
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
	var value interface{}
	parameters := make(map[string]string)
	for i := 1; i < len(actionArray); i++ {
		if i == 1 && !(strings.Contains(actionArray[i], "=")) {
			value = actionArray[i]
		} else {
			param := strings.Split(actionArray[i], "=")
			paramName := param[0]
			paramValue := param[1]
			parameters[paramName] = paramValue
		}
	}
	return &Action{
		ExposedThing:   exposedThing,
		Ref:            ref,
		AutomationName: automationName,
		ExposedAction:  exposedAction,
		Value:          value,
		Parameters:     parameters,
	}, nil
}

func (a *Action) run() error {
	// TODO: use parameters
	if a.ExposedAction == nil {
		return errors.New("missing action handler")
	}
	_, err := a.ExposedAction.Run(a.ExposedThing, a.Ref, a.Value, a.Parameters)
	// Note we don't do anything with the output
	return err
}
