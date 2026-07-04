package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var moving, _ = dataSchema.NewBoolean()
var Moving = Meta{
	Title:       "Moving",
	Description: "Whether the device is currently moving",
	Data:        moving,
}
