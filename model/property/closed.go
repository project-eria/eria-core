package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var closed, _ = dataSchema.NewBoolean()
var Closed = Meta{
	Title:       "Closed",
	Description: "Whether the device is fully closed (end-of-travel contact)",
	Data:        closed,
}
