package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var strength, _ = dataSchema.NewNumber(
	dataSchema.NumberUnit("m/s"),
)
var WindStrength = Meta{
	Title:       "Wind Strength",
	Description: "Strength of the wind",
	Data:        strength,
}
