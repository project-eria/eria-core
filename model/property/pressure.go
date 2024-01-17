package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var pressure, _ = dataSchema.NewInteger(
	dataSchema.IntegerUnit("hPa"),
)
var Pressure = Meta{
	Title:       "Pressure",
	Description: "A barometric pressure meter",
	Data:        pressure,
}
