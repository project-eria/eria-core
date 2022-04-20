package model

type SchemaAction struct {
	Id   string
	Name string
}

type ActionDesc struct {
	Title       string
	Description string
}

var CapabilitiesActions = map[string][]SchemaAction{
	"LightBasic": []SchemaAction{
		{Id: "toggle", Name: "ToggleAction"},
	},
	"LightDimmer": []SchemaAction{
		{Id: "fade", Name: "FadeAction"},
	},
	"ShutterBasic": []SchemaAction{
		{Id: "open", Name: "OpenAction"},
		{Id: "close", Name: "CloseAction"},
		{Id: "stop", Name: "StopAction"},
	},
	"ShutterPosition": []SchemaAction{
		{Id: "setPosition", Name: "SetPositionAction"},
	},
}

// actions schemas list
var Actions = map[string]ActionDesc{
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
