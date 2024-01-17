package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var position, _ = dataSchema.NewInteger(
	dataSchema.IntegerDefault(0),
	dataSchema.IntegerUnit("%"),
	dataSchema.IntegerMin(0),
	dataSchema.IntegerMax(100),
)

var Position = Meta{
	Title:       "Position",
	Description: "The position property from 0-100",
	Data:        position,
}
