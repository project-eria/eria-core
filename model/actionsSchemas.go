package model

import (
	"github.com/project-eria/go-wot/dataSchema"
)

type SchemaAction struct {
	Id   string
	Name string
}

type ActionDesc struct {
	Title       string
	Description string
	Input       dataSchema.Data
	Output      dataSchema.Data
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
	"ORPSensor": []SchemaAction{
		{Id: "calibrateORP", Name: "CalibrateORPAction"},
	},
	"PHSensor": []SchemaAction{
		{Id: "calibratePH", Name: "CalibratePHAction"},
	},
	"TemperatureSensor": []SchemaAction{
		{Id: "calibrateTemperature", Name: "CalibrateTemperatureAction"},
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
	"CalibrateORPAction": {
		Title:       "Calibrate ORP",
		Description: "Calibrate the ORP value with buffer solution",
		Input:       dataSchema.NewNumber(0.0, "mV", -2000, 2000),
	},
	"CalibratePHAction": {
		Title:       "Calibrate PH",
		Description: "Calibrate the PH value with buffer solution",
		Input:       dataSchema.NewNumber(0.0, "", 0, 14),
	},
	"CalibrateTemperatureAction": {
		Title:       "Calibrate Temperature",
		Description: "Calibrate the Temperature value",
		Input:       dataSchema.NewNumber(0.0, "°C", -1000, 1000),
	},
}
