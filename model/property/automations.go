package propertyModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var automations, _ = dataSchema.NewArray()

// Automations is the read-only property exposing the list of currently
// scheduled automations (device, action, next run date/time).
var Automations = Meta{
	Title:       "Scheduled automations",
	Description: "List of the currently scheduled automations, with their devices, action and next run date/time",
	Data:        automations,
	Options: []interaction.PropertyOption{
		interaction.PropertyReadOnly(true),
		interaction.PropertyObservable(false), // read on demand, no change stream
	},
}
