package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var uv, _ = dataSchema.NewNumber(
	dataSchema.NumberMin(0),
	dataSchema.NumberMax(20),
)
var UVIndex = Meta{
	Title:       "UV Index",
	Description: "Strength of the sunburn-producing ultraviolet (UV) radiation",
	Data:        uv,
}
