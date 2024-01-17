package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var orp, _ = dataSchema.NewNumber(
	dataSchema.NumberUnit("mV"),
)
var ORP = Meta{
	Title:       "ORP/Redox",
	Description: "A reduction / oxidation potential meter",
	Data:        orp,
}
