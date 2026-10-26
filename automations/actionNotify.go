package automations

import (
	"errors"
	"fmt"
	"strings"

	zlog "github.com/rs/zerolog/log"
)

const notifyPrefix = "notify|"

// NoticesThingTag is the tag of the remote `service-notices` Thing, looked up
// through the eria-core consumer (declare it in `remoteThings` with this tag).
var NoticesThingTag = "service-notices"

// actionInvoker is the part of a consumed Thing able to invoke a WoT action
// (implemented by *eriaconsumer.EriaThing).
type actionInvoker interface {
	InvokeAction(name string, uriVariables map[string]interface{}, value interface{}) (interface{}, error)
}

// actionResolver is implemented by actions that must be undone when the
// automation's conditions go back to false (see Automation.plan).
type actionResolver interface {
	resolve() error
}

// NotifyAction raises a notice on the `service-notices` Thing, and resolves it
// when the automation's conditions stop matching.
//
// Syntax: `notify|key=<key>|severity=<info|warning|alert>|title=<title>|body=<body>|audience=<all|profile>`
//   - key and title are required; key dedupes the notice
//   - severity defaults to info, audience to all
type NotifyAction struct {
	AutomationName string
	Key            string
	Severity       string
	Title          string
	Body           string
	Audience       string
	Thing          string // first thing of the automation, if any
}

func getNotifyAction(thingRefs []string, automationName string, actionStr string) (*NotifyAction, error) {
	action := &NotifyAction{
		AutomationName: automationName,
		Severity:       "info",
		Audience:       "all",
	}
	if len(thingRefs) > 0 {
		action.Thing = thingRefs[0]
	}
	actionArray := strings.Split(actionStr, "|")
	for i := 1; i < len(actionArray); i++ {
		// SplitN, not Split: the value may itself contain '='
		param := strings.SplitN(actionArray[i], "=", 2)
		if len(param) != 2 {
			return nil, fmt.Errorf("invalid notify parameter %q", actionArray[i])
		}
		value := strings.TrimSpace(param[1])
		switch param[0] {
		case "key":
			action.Key = value
		case "title":
			action.Title = value
		case "body":
			action.Body = value
		case "severity":
			if value != "info" && value != "warning" && value != "alert" {
				return nil, fmt.Errorf("unsupported notify severity %q", value)
			}
			action.Severity = value
		case "audience":
			if value != "all" && value != "profile" {
				return nil, fmt.Errorf("unsupported notify audience %q", value)
			}
			action.Audience = value
		default:
			return nil, fmt.Errorf("unknown notify parameter %q", param[0])
		}
	}
	if action.Key == "" {
		return nil, errors.New("missing notify key")
	}
	if action.Title == "" {
		return nil, errors.New("missing notify title")
	}
	return action, nil
}

func (a *NotifyAction) run() error {
	input := map[string]interface{}{
		"key":      a.Key,
		"source":   "house",
		"severity": a.Severity,
		"title":    a.Title,
		"audience": a.Audience,
	}
	if a.Body != "" {
		input["body"] = a.Body
	}
	if a.Thing != "" {
		input["thing"] = a.Thing
	}
	return a.invoke("raise", input)
}

func (a *NotifyAction) resolve() error {
	return a.invoke("resolve", map[string]interface{}{"key": a.Key})
}

func (a *NotifyAction) invoke(name string, input map[string]interface{}) error {
	var thing actionInvoker
	if _consumer != nil {
		// The interface holds a typed nil when the tag is unknown: check both
		if t, ok := _consumer.ThingFromTag(NoticesThingTag).(actionInvoker); ok && t != nil {
			thing = t
		}
	}
	if thing == nil {
		err := errors.New("notices thing not found")
		zlog.Error().Err(err).Str("automation", a.AutomationName).Str("tag", NoticesThingTag).Msg("[automations:NotifyAction] Can't " + name)
		return err
	}
	if _, err := thing.InvokeAction(name, nil, input); err != nil {
		zlog.Error().Err(err).Str("automation", a.AutomationName).Str("key", a.Key).Msg("[automations:NotifyAction] Can't " + name)
		return err
	}
	zlog.Info().Str("automation", a.AutomationName).Str("key", a.Key).Msg("[automations:NotifyAction] Notice " + name)
	return nil
}
