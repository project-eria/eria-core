package propertyModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var position, _ = dataSchema.NewInteger(
	dataSchema.IntegerDefault(0),
	dataSchema.IntegerUnit("%"),
	dataSchema.IntegerMin(0),
	dataSchema.IntegerMax(100),
)

var Position = Meta{
	Title:       "Position",
	Description: "The position property from 0-100",
	Data:        position,
	// State reported by the device; change it through the actions
	Options: []interaction.PropertyOption{interaction.PropertyReadOnly(true)},
}
