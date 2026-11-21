package eventModel

// Alarm is the start of a detection; its data names what was detected
// (MotionAlarm, PeopleDetect, VehicleDetect, DogCatDetect…).
var Alarm = Meta{
	Title:       "Alarm",
	Description: "Start of a detection; data is what was detected (MotionAlarm, PeopleDetect, VehicleDetect, DogCatDetect…)",
}
