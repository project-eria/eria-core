package model

import (
	"errors"

	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// stringInput := dataSchema.NewString(false)
// stringOutput := dataSchema.NewString(false)

// AddActionFromSchema return an action from schema @type
func AddActionFromSchema(t *thing.Thing, id string, schema string) (*interaction.Action, error) {
	log.Info().Str("schema", schema).Msg("[thing:AddActionFromSchema] Adding action")

	if meta, in := actions[schema]; in {
		// var input dataSchema.Data
		// var output dataSchema.Data
		action := interaction.NewAction(
			id,
			meta.Title,
			meta.Description,
			nil,
			nil,
		)
		// TODO (remove?) action.ATtype = schema
		t.AddAction(action)
		return action, nil
	}
	return nil, errors.New("Action schema '" + schema + "' not found")
}

// actions schemas list
var actions = map[string]struct {
	Title       string
	Description string
}{
	"ToggleAction": {
		Title:       "Toggle",
		Description: "Toggles a boolean state on and off",
	},
	"FadeAction": {
		Title:       "Fade",
		Description: "Fade the lamp to a given level",
		// TODO
		// Input: thing.ActionInputMeta{
		// 	Type: "object",
		// 	Properties: map[string]thing.ActionInputPropertyMeta{
		// 		"level": thing.ActionInputPropertyMeta{
		// 			Type:    "integer",
		// 			Minimum: 0,
		// 			Maximum: 100,
		// 			Unit:    "%",
		// 		},
		// 		"duration": thing.ActionInputPropertyMeta{
		// 			Type:    "integer",
		// 			Minimum: 0,
		// 			Unit:    "ms",
		// 		},
		// 	},
		// },
	},
	"OpenAction": {
		Title:       "Open",
		Description: "Open a motorised device",
	},
	"CloseAction": {
		Title:       "Close",
		Description: "Close a motorised device",
	},
	"StopAction": {
		Title:       "Stop",
		Description: "Stop an openning/closing action on a motorised device",
	},
	"SetPositionAction": {
		Title:       "Set Position",
		Description: "Set a particular position on an motorized device",
		// TODO
		// Input: thing.ActionInputMeta{
		// 	Type: "object",
		// 	Properties: map[string]thing.ActionInputPropertyMeta{
		// 		"target": thing.ActionInputPropertyMeta{
		// 			Type:    "integer",
		// 			Minimum: 0,
		// 			Maximum: 100,
		// 			Unit:    "%",
		// 		},
		// 	},
		// },
	},
}
