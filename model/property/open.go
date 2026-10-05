package propertyModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var open, _ = dataSchema.NewBoolean()
var Open = Meta{
	Title:       "Open",
	Description: "Whether the device is open",
	Data:        open,
	// State reported by the device; change it through the actions
	Options: []interaction.PropertyOption{interaction.PropertyReadOnly(true)},
}
