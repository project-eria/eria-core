package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var angle, _ = dataSchema.NewInteger(
	dataSchema.IntegerUnit("°"),
	dataSchema.IntegerMin(0),
	dataSchema.IntegerMax(360),
)
var WindAngle = Meta{
	Title:       "Wind Angle",
	Description: "Direction of the wind",
	Data:        angle,
}
