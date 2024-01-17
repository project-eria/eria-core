package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var open, _ = dataSchema.NewBoolean()
var Open = Meta{
	Title:       "Open",
	Description: "Whether the device is open",
	Data:        open,
}
