package automations

import (
	"errors"
	"strings"

	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

type ActionRunner interface {
	run() error
}

type Action struct {
	AutomationName string
	ExposedThings  map[string]producer.ExposedThing
	Ref            string
	Value          interface{}
	Parameters     map[string]string
}

/**
 * `<action>|<value>|<param name>=<value>|<param name>=<value>`
 */
func getAction(thingRefs []string, automationName string, actionStr string) (*Action, error) {
	actionStr = strings.TrimSpace(actionStr)
	if actionStr == "" {
		return nil, errors.New("missing action configuration") // Skip this automation
	}
	actionArray := strings.Split(actionStr, "|")
	ref := actionArray[0]

	things := map[string]producer.ExposedThing{}
	oneValid := false
	for _, thingRef := range thingRefs {
		// Get the thing
		if exposedThing, ok := _exposedThings[thingRef]; ok && exposedThing != nil {
			// Check if the action exists
			_, err := exposedThing.ExposedAction(ref)
			if err != nil {
				zlog.Error().Err(err).Str("automation", automationName).Str("thing", thingRef).Str("action", ref).Msg("[automations:getAction] Action not found")
				continue // jump to next thing
			}
			oneValid = true
			things[thingRef] = exposedThing
		} else {
			zlog.Error().Str("automation", automationName).Str("thing", thingRef).Str("action", ref).Msg("[automations:getAction] Thing not found")
		}
	}

	if !oneValid {
		return nil, errors.New("requested things action not found")
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
		ExposedThings:  things,
		Ref:            ref,
		AutomationName: automationName,
		Value:          value,
		Parameters:     parameters,
	}, nil
}

func (a *Action) run() error {
	// TODO: use parameters
	oneHasRun := false
	for key, thing := range a.ExposedThings {
		exposedAction, _ := thing.ExposedAction(a.Ref)
		if exposedAction == nil {
			return errors.New("missing action handler")
		}
		_, err := exposedAction.Run(thing, a.Ref, a.Value, a.Parameters)
		// Note we don't do anything with the output
		if err != nil {
			zlog.Error().Err(err).Str("automation", a.AutomationName).Str("thing", key).Str("action", a.Ref).Msg("[automations:getAction] Action failed")
			continue
		}
		oneHasRun = true
	}
	if !oneHasRun {
		return errors.New("all things failed to run action")
	}
	return nil
}
