package propertyModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var closed, _ = dataSchema.NewBoolean()
var Closed = Meta{
	Title:       "Closed",
	Description: "Whether the device is fully closed (end-of-travel contact)",
	Data:        closed,
	// State reported by the device; change it through the actions
	Options: []interaction.PropertyOption{interaction.PropertyReadOnly(true)},
}
