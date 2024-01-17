package actionModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

var inputTemperature, _ = dataSchema.NewNumber(
	dataSchema.NumberDefault(0.0),
	dataSchema.NumberUnit("°C"),
	dataSchema.NumberMin(-1000),
	dataSchema.NumberMax(1000),
)

var CalibrateTemperature = Meta{
	Title:       "Calibrate Temperature",
	Description: "Calibrate the Temperature value",
	Options: []interaction.ActionOption{
		interaction.ActionInput(&inputTemperature),
	},
}
