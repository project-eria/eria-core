package model

import (
	actionModel "github.com/project-eria/eria-core/model/action"
	propertyModel "github.com/project-eria/eria-core/model/property"
)

var Models = map[string]Model{
	"EriaAppController": {
		Version: "1.0.0",
		Properties: map[string]ModelProperty{
			"logLevel": {Meta: propertyModel.LogLevel},
		},
	},
	"LightBasic": {
		Properties: map[string]ModelProperty{
			"on": {DefaultValue: false, Meta: propertyModel.OnOff},
		},
		Actions: map[string]ModelAction{
			"toggle": {Meta: actionModel.Toggle},
		},
	},
	"LightDimmer": {
		Properties: map[string]ModelProperty{
			"on":         {DefaultValue: false, Meta: propertyModel.OnOff},
			"brightness": {DefaultValue: 0, Meta: propertyModel.Brightness},
		},
		Actions: map[string]ModelAction{
			"fade":   {Meta: actionModel.Fade},
			"toggle": {Meta: actionModel.Toggle},
		},
	},
	"ShutterBasic": {
		Properties: map[string]ModelProperty{
			"open": {DefaultValue: false, Meta: propertyModel.Open},
		},
		Actions: map[string]ModelAction{
			"open":  {Meta: actionModel.Open},
			"close": {Meta: actionModel.Close},
			"stop":  {Meta: actionModel.Stop},
		},
	},
	"ShutterPosition": {
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
	"TemperatureSensor": {
		Properties: map[string]ModelProperty{
			"temperature": {DefaultValue: 0.0, Meta: propertyModel.Temperature},
		},
		Actions: map[string]ModelAction{
			"calibrateTemperature": {Meta: actionModel.CalibrateTemperature},
		},
	},
	"VoltageSensor": {
		Properties: map[string]ModelProperty{
			"volts": {DefaultValue: 0.0, Meta: propertyModel.Voltage},
		},
	},
	"HygrometerSensor": {
		Properties: map[string]ModelProperty{
			"humidity": {DefaultValue: 0, Meta: propertyModel.Humidity},
		},
	},
	"BarometerSensor": {
		Properties: map[string]ModelProperty{
			"pressure": {DefaultValue: 0, Meta: propertyModel.Pressure},
		},
	},
	"WindgaugeSensor": {
		Properties: map[string]ModelProperty{
			"windStrength": {DefaultValue: 0.0, Meta: propertyModel.WindStrength},
			"windAngle":    {DefaultValue: 0, Meta: propertyModel.WindAngle},
		},
	},
	"WaterMeter": {
		Properties: map[string]ModelProperty{
			"liters": {DefaultValue: 0.0, Meta: propertyModel.Liters},
		},
	},
	"UVSensor": {
		Properties: map[string]ModelProperty{
			"uvIndex": {DefaultValue: 0.0, Meta: propertyModel.UVIndex},
		},
	},
	"PHSensor": {
		Properties: map[string]ModelProperty{
			"ph": {DefaultValue: 0.0, Meta: propertyModel.PH},
		},
		Actions: map[string]ModelAction{
			"calibratePH": {Meta: actionModel.CalibratePH},
		},
	},
	"ORPSensor": {
		Properties: map[string]ModelProperty{
			"orp": {DefaultValue: 0, Meta: propertyModel.ORP},
		},
		Actions: map[string]ModelAction{
			"calibrateORP": {Meta: actionModel.CalibrateORP},
		},
	},
}
