package eria

import (
	"errors"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// AddPropertyFromSchema return an property from schema @type
func AddPropertyFromSchema(t *thing.Thing, id string, defaultValue interface{}, schema string) (*interaction.Property, error) {
	log.Info().Str("schema", schema).Msg("[thing:AddPropertyFromSchema] Adding property")

	if meta, in := properties[schema]; in {
		var data dataSchema.Data
		switch meta.Type {
		case "boolean":
			data = dataSchema.NewBoolean(defaultValue.(bool))
		case "integer":
			data = dataSchema.NewInteger(defaultValue.(int), meta.Unit, meta.Minimum, meta.Maximum)
		case "number":
			data = dataSchema.NewNumber(defaultValue.(float64), meta.Unit, meta.Minimum, meta.Maximum)
		}
		property := interaction.NewProperty(
			id,
			meta.Title,
			meta.Description,
			false,
			false,
			true,
			data,
		)
		t.AddProperty(property)
		// TODO (remove?) property.ATtype = schema
		return property, nil
	}
	return nil, errors.New("Property schema '" + schema + "' not found")
}

var properties = map[string]struct {
	Title       string
	Description string
	Type        string
	ReadOnly    bool
	Unit        string
	Minimum     uint16
	Maximum     uint16
}{
	"OnOffProperty": {
		Title:       "On/Off",
		Description: "Whether the device is turned on",
		Type:        "boolean",
	},
	"BrightnessProperty": {
		Title:       "Brightness",
		Description: "The brightness level from 0-100",
		Type:        "integer",
		Unit:        "%",
		Minimum:     0,
		Maximum:     100,
	},
	"OpenProperty": {
		Title:       "Open",
		Description: "Whether the device is open",
		Type:        "boolean",
	},
	"PositionProperty": {
		Title:       "Position",
		Description: "The position property from 0-100",
		Type:        "integer",
		Unit:        "%",
		Minimum:     0,
		Maximum:     100,
	},
	"TemperatureProperty": {
		Title:       "Temperature",
		Description: "An ambient Celsius temperature sensor",
		Type:        "number",
		Unit:        "°C",
	},
	"VoltageProperty": {
		Title:       "Voltage",
		Description: "An voltage value",
		Type:        "number",
		Unit:        "V",
	},
	"PhProperty": {
		Title:       "pH",
		Description: "A pH meter to mesure acidity/basicity",
		Type:        "number",
	},
	"OrpProperty": {
		Title:       "ORP/Redox",
		Description: "A reduction / oxidation potential meter",
		Type:        "integer",
		Unit:        "mV",
	},
	"HumidityProperty": {
		Title:       "Humidity",
		Description: "The humidity property from 0-100",
		Type:        "integer",
		Unit:        "%",
		Minimum:     0,
		Maximum:     100,
	},
	"PressureProperty": {
		Title:       "Pressure",
		Description: "A barometric pressure meter",
		Type:        "integer",
		Unit:        "hPa",
	},
	"WindAngleProperty": {
		Title:       "Wind Angle",
		Description: "Direction of the wind",
		Type:        "integer",
		Unit:        "°",
		Minimum:     0,
		Maximum:     360,
	},
	"WindStrengthProperty": {
		Title:       "Wind Strength",
		Description: "Strength of the wind",
		Type:        "number",
		Unit:        "m/s",
	},
	"LitersProperty": {
		Title:       "Liters Volume",
		Description: "Number of liters",
		Type:        "number",
		Unit:        "l",
	},
}
