package model

type SchemaEvent struct {
	Id   string
	Name string
}

type EventDesc struct {
	Title       string
	Description string
}

var CapabilitiesEvents = map[string][]SchemaEvent{}

// events schemas list
var Events = map[string]EventDesc{}
