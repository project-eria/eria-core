package model

import (
	actionModel "github.com/project-eria/eria-core/model/action"
	propertyModel "github.com/project-eria/eria-core/model/property"
)

var Models = map[string]Model{
	"EriaAppController": Model{
		Version: "1.0.0",
		Properties: map[string]ModelProperty{
			"logLevel": {Meta: propertyModel.LogLevel},
		},
	},
	"LightBasic": Model{
		Properties: map[string]ModelProperty{
			"on": {DefaultValue: false, Meta: propertyModel.OnOff},
		},
		Actions: map[string]ModelAction{
			"toggle": {Meta: actionModel.Toggle},
		},
	},
	"LightDimmer": Model{
		Properties: map[string]ModelProperty{
			"on":         {DefaultValue: false, Meta: propertyModel.OnOff},
			"brightness": {DefaultValue: 0, Meta: propertyModel.Brightness},
		},
		Actions: map[string]ModelAction{
			"fade":   {Meta: actionModel.Fade},
			"toggle": {Meta: actionModel.Toggle},
		},
	},
	"ShutterBasic": Model{
		Properties: map[string]ModelProperty{
			"open": {DefaultValue: false, Meta: propertyModel.Open},
		},
		Actions: map[string]ModelAction{
			"open":  {Meta: actionModel.Open},
			"close": {Meta: actionModel.Close},
			"stop":  {Meta: actionModel.Stop},
		},
	},
	"ShutterPosition": Model{
		Properties: map[string]ModelProperty{
			"open":     {DefaultValue: false, Meta: propertyModel.Open},
			"position": {DefaultValue: 0, Meta: propertyModel.Position},
		},
		Actions: map[string]ModelAction{
			"open":        {Meta: actionModel.Open},
			"close":       {Meta: actionModel.Close},
			"stop":        {Meta: actionModel.Stop},
			"setPosition": {Meta: actionModel.SetPosition},
		},
	},
	"TemperatureSensor": Model{
		Properties: map[string]ModelProperty{
			"temperature": {DefaultValue: 0.0, Meta: propertyModel.Temperature},
		},
		Actions: map[string]ModelAction{
			"calibrateTemperature": {Meta: actionModel.CalibrateTemperature},
		},
	},
	"VoltageSensor": Model{
		Properties: map[string]ModelProperty{
			"volts": {DefaultValue: 0.0, Meta: propertyModel.Voltage},
		},
	},
	"HygrometerSensor": Model{
		Properties: map[string]ModelProperty{
			"humidity": {DefaultValue: 0, Meta: propertyModel.Humidity},
		},
	},
	"BarometerSensor": Model{
		Properties: map[string]ModelProperty{
			"pressure": {DefaultValue: 0, Meta: propertyModel.Pressure},
		},
	},
	"WindgaugeSensor": Model{
		Properties: map[string]ModelProperty{
			"windStrength": {DefaultValue: 0.0, Meta: propertyModel.WindStrength},
			"windAngle":    {DefaultValue: 0, Meta: propertyModel.WindAngle},
		},
	},
	"WaterMeter": Model{
		Properties: map[string]ModelProperty{
			"liters": {DefaultValue: 0.0, Meta: propertyModel.Liters},
		},
	},
	"UVSensor": Model{
		Properties: map[string]ModelProperty{
			"uvIndex": {DefaultValue: 0.0, Meta: propertyModel.UVIndex},
		},
	},
	"PHSensor": Model{
		Properties: map[string]ModelProperty{
			"ph": {DefaultValue: 0.0, Meta: propertyModel.PH},
		},
		Actions: map[string]ModelAction{
			"calibratePH": {Meta: actionModel.CalibratePH},
		},
	},
	"ORPSensor": Model{
		Properties: map[string]ModelProperty{
			"orp": {DefaultValue: 0, Meta: propertyModel.ORP},
		},
		Actions: map[string]ModelAction{
			"calibrateORP": {Meta: actionModel.CalibrateORP},
		},
	},
}
