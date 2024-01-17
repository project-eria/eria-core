package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var temperature, _ = dataSchema.NewNumber(
	dataSchema.NumberUnit("°C"),
)
var Temperature = Meta{
	Title:       "Temperature",
	Description: "An ambient Celsius temperature sensor",
	Data:        temperature,
}
