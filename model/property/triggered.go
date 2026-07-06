package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var triggered, _ = dataSchema.NewBoolean()
var Triggered = Meta{
	Title:       "Triggered",
	Description: "Whether the trigger is currently active",
	Data:        triggered,
}
