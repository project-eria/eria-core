package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var humidity, _ = dataSchema.NewInteger(
	dataSchema.IntegerDefault(0),
	dataSchema.IntegerUnit("%"),
	dataSchema.IntegerMin(0),
	dataSchema.IntegerMax(100),
)

var Humidity = Meta{
	Title:       "Humidity",
	Description: "The humidity property from 0-100",
	Data:        humidity,
}
