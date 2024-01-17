package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var ph, _ = dataSchema.NewNumber()
var PH = Meta{
	Title:       "pH",
	Description: "A pH meter to mesure acidity/basicity",
	Data:        ph,
}
