package propertyModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var triggered, _ = dataSchema.NewBoolean()
var Triggered = Meta{
	Title:       "Triggered",
	Description: "Whether the trigger is currently active",
	Data:        triggered,
	// State reported by the device; change it through the actions
	Options: []interaction.PropertyOption{interaction.PropertyReadOnly(true)},
}
