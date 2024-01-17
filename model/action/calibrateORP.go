package actionModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var inputORP, _ = dataSchema.NewNumber(
	dataSchema.NumberDefault(0.0),
	dataSchema.NumberUnit("mV"),
	dataSchema.NumberMin(-2000),
	dataSchema.NumberMax(2000),
)

var CalibrateORP = Meta{
	Title:       "Calibrate ORP",
	Description: "Calibrate the ORP value with buffer solution",
	Options: []interaction.ActionOption{
		interaction.ActionInput(&inputORP),
	},
}
