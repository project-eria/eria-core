package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var voltage, _ = dataSchema.NewNumber(
	dataSchema.NumberUnit("V"),
)
var Voltage = Meta{
	Title:       "Voltage",
	Description: "An voltage value",
	Data:        voltage,
}
