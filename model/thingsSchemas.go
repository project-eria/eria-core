package model

import (
	"errors"

	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// NewFromSchemas return a thing from schemas @type
func NewFromSchemas(urn string, title string, description string, capabilities []string) (*thing.Thing, error) {
	t, err := thing.New(urn, title, description, nil)
	if err != nil {
		return nil, err
	}

	if err := AddSchemas(t, capabilities, ""); err != nil {
		return nil, err
	}

	return t, nil
}

// AddSchemas add capabilities to a thing using schemas @type
func AddSchemas(t *thing.Thing, capabilities []string, prefix string) error {
	// TODO (remove?) t.SetContext("http://www.w3.org/ns/td")

	for _, capability := range capabilities {
		if err := AddSchema(t, capability, prefix); err != nil {
			return err
		}
	}

	return nil
}

// AddSchema add a capability to a thing using schema @type
func AddSchema(t *thing.Thing, capability string, prefix string) error {
	log.Info().Str("capability", capability).Msg("[thing:AddSchema] Adding capability")
	propertiesSchemas, inProperties := capabilitiesProperties[capability]
	if inProperties {
		for _, schema := range propertiesSchemas {
			id := prefix + schema.id
			if _, err := AddPropertyFromSchema(t, id, schema.defaultValue, schema.name); err != nil {
				return err
			}
		}
	}
	actionsSchemas, inActions := capabilitiesActions[capability]
	if inActions {
		for _, schema := range actionsSchemas {
			id := prefix + schema.id
			if _, err := AddActionFromSchema(t, id, schema.name); err != nil {
				return err
			}
		}
	}
	if !inProperties && !inActions {
		return errors.New("Capability schema '" + capability + "' not found")
	}

	// TODO (remove?) t.AddType(capability)

	return nil
}

type schemaProperty struct {
	id           string
	defaultValue interface{}
	name         string
}

type schemaAction struct {
	id   string
	name string
}

var capabilitiesProperties = map[string][]schemaProperty{
	"LightBasic": []schemaProperty{
		{id: "on", defaultValue: false, name: "OnOffProperty"},
	},
	"LightDimmer": []schemaProperty{
		{id: "brightness", defaultValue: 0, name: "BrightnessProperty"},
	},
	"ShutterBasic": []schemaProperty{
		{id: "open", defaultValue: false, name: "OpenProperty"},
	},
	"ShutterPosition": []schemaProperty{
		{id: "position", defaultValue: 0, name: "PositionProperty"},
	},
	"TemperatureSensor": []schemaProperty{
		{id: "temperature", defaultValue: 0.0, name: "TemperatureProperty"},
	},
	"VoltageSensor": []schemaProperty{
		{id: "volts", defaultValue: 0.0, name: "VoltageProperty"},
	},
	"PoolMonitor": []schemaProperty{
		{id: "temperature", defaultValue: 0.0, name: "TemperatureProperty"},
		{id: "ph", defaultValue: 0.0, name: "PhProperty"},
		{id: "orp", defaultValue: 0, name: "OrpProperty"},
	},
	"HygrometerSensor": []schemaProperty{
		{id: "humidity", defaultValue: 0, name: "HumidityProperty"},
	},
	"BarometerSensor": []schemaProperty{
		{id: "pressure", defaultValue: 0, name: "PressureProperty"},
	},
	"WindgaugeSensor": []schemaProperty{
		{id: "windStrength", defaultValue: 0.0, name: "WindStrengthProperty"},
		{id: "windAngle", defaultValue: 0, name: "WindAngleProperty"},
	},
	"WaterMeter": []schemaProperty{
		{id: "liters", defaultValue: 0.0, name: "LitersProperty"},
	},
}
var capabilitiesActions = map[string][]schemaAction{
	"LightBasic": []schemaAction{
		{id: "toggle", name: "ToggleAction"},
	},
	"LightDimmer": []schemaAction{
		{id: "fade", name: "FadeAction"},
	},
	"ShutterBasic": []schemaAction{
		{id: "open", name: "OpenAction"},
		{id: "close", name: "CloseAction"},
		{id: "stop", name: "StopAction"},
	},
	"ShutterPosition": []schemaAction{
		{id: "setPosition", name: "SetPositionAction"},
	},
}
