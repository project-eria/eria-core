package model

type SchemaProperty struct {
	Id           string
	DefaultValue interface{}
	Name         string
}

type PropertyDesc struct {
	Title       string
	Description string
	Type        string
	ReadOnly    bool
	Unit        string
	Minimum     int16
	Maximum     int16
	MinLength   uint16
	MaxLength   uint16
	Pattern     string
}

var CapabilitiesProperties = map[string][]SchemaProperty{
	"LightBasic": []SchemaProperty{
		{Id: "on", DefaultValue: false, Name: "OnOffProperty"},
	},
	"LightDimmer": []SchemaProperty{
		{Id: "brightness", DefaultValue: 0, Name: "BrightnessProperty"},
	},
	"ShutterBasic": []SchemaProperty{
		{Id: "open", DefaultValue: false, Name: "OpenProperty"},
	},
	"ShutterPosition": []SchemaProperty{
		{Id: "position", DefaultValue: 0, Name: "PositionProperty"},
	},
	"TemperatureSensor": []SchemaProperty{
		{Id: "temperature", DefaultValue: 0.0, Name: "TemperatureProperty"},
	},
	"VoltageSensor": []SchemaProperty{
		{Id: "volts", DefaultValue: 0.0, Name: "VoltageProperty"},
	},
	"HygrometerSensor": []SchemaProperty{
		{Id: "humidity", DefaultValue: 0, Name: "HumidityProperty"},
	},
	"BarometerSensor": []SchemaProperty{
		{Id: "pressure", DefaultValue: 0, Name: "PressureProperty"},
	},
	"WindgaugeSensor": []SchemaProperty{
		{Id: "windStrength", DefaultValue: 0.0, Name: "WindStrengthProperty"},
		{Id: "windAngle", DefaultValue: 0, Name: "WindAngleProperty"},
	},
	"WaterMeter": []SchemaProperty{
		{Id: "liters", DefaultValue: 0.0, Name: "LitersProperty"},
	},
	"UVSensor": []SchemaProperty{
		{Id: "uvIndex", DefaultValue: 0.0, Name: "UVIndexProperty"},
	},
	"PHSensor": []SchemaProperty{
		{Id: "ph", DefaultValue: 0.0, Name: "PHProperty"},
	},
	"ORPSensor": []SchemaProperty{
		{Id: "orp", DefaultValue: 0, Name: "ORPProperty"},
	},
}

var Properties = map[string]PropertyDesc{
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
	"PHProperty": {
		Title:       "pH",
		Description: "A pH meter to mesure acidity/basicity",
		Type:        "number",
	},
	"ORPProperty": {
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
	"UVIndexProperty": {
		Title:       "UV Index",
		Description: "Strength of the sunburn-producing ultraviolet (UV) radiation",
		Type:        "number",
		Minimum:     0,
		Maximum:     20,
	},
}
