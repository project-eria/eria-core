package actionModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var inputPosition, _ = dataSchema.NewInteger(
	dataSchema.IntegerDefault(0),
	dataSchema.IntegerUnit("%"),
	dataSchema.IntegerMin(0),
	dataSchema.IntegerMax(100),
)

var SetPosition = Meta{
	Title:       "Set Position",
	Description: "Set a particular position on an motorized device",
	Options: []interaction.ActionOption{
		interaction.ActionInput(&inputPosition),
	},
}
