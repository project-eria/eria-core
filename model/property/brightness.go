package propertyModel

import (
	"github.com/project-eria/go-wot/dataSchema"
)

var brightness, _ = dataSchema.NewInteger(
	dataSchema.IntegerDefault(0),
	dataSchema.IntegerUnit("%"),
	dataSchema.IntegerMin(0),
	dataSchema.IntegerMax(100),
)
var Brightness = Meta{
	Title:       "Brightness",
	Description: "The brightness level from 0-100",
	Data:        brightness,
}
