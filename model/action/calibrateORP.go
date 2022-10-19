package actionModel

import "github.com/project-eria/go-wot/dataSchema"

var CalibrateORP = Meta{
	Title:       "Calibrate ORP",
	Description: "Calibrate the ORP value with buffer solution",
	Input:       dataSchema.NewNumber(0.0, "mV", -2000, 2000),
}
