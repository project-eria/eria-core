package actionModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var inputPH, _ = dataSchema.NewNumber(
	dataSchema.NumberDefault(0.0),
	dataSchema.NumberMin(0),
	dataSchema.NumberMax(14),
)

var CalibratePH = Meta{
	Title:       "Calibrate PH",
	Description: "Calibrate the PH value with buffer solution",
	Options: []interaction.ActionOption{
		interaction.ActionInput(&inputPH),
	},
}
