package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var liters, _ = dataSchema.NewNumber(
	dataSchema.NumberUnit("l"),
)
var Liters = Meta{
	Title:       "Liters Volume",
	Description: "Number of liters",
	Data:        liters,
}
