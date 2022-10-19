package actionModel

import "github.com/project-eria/go-wot/dataSchema"

var CalibratePH = Meta{
	Title:       "Calibrate PH",
	Description: "Calibrate the PH value with buffer solution",
	Input:       dataSchema.NewNumber(0.0, "", 0, 14),
}
