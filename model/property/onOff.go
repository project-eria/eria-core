package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var onOff, _ = dataSchema.NewBoolean()
var OnOff = Meta{
	Title:       "On/Off",
	Description: "Whether the device is turned on",
	Data:        onOff,
}
