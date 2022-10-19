package actionModel

import "github.com/project-eria/go-wot/dataSchema"

var CalibrateTemperature = Meta{
	Title:       "Calibrate Temperature",
	Description: "Calibrate the Temperature value",
	Input:       dataSchema.NewNumber(0.0, "°C", -1000, 1000),
}
