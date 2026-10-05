package propertyModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var moving, _ = dataSchema.NewBoolean()
var Moving = Meta{
	Title:       "Moving",
	Description: "Whether the device is currently moving",
	Data:        moving,
	// State reported by the device; change it through the actions
	Options: []interaction.PropertyOption{interaction.PropertyReadOnly(true)},
}
