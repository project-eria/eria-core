package propertyModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var onOff, _ = dataSchema.NewBoolean()
var OnOff = Meta{
	Title:       "On/Off",
	Description: "Whether the device is turned on",
	Data:        onOff,
	// State reported by the device; change it through the actions
	Options: []interaction.PropertyOption{interaction.PropertyReadOnly(true)},
}
