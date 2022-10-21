package model

import (
	actionModel "github.com/project-eria/eria-core/model/action"
	eventModel "github.com/project-eria/eria-core/model/event"
	propertyModel "github.com/project-eria/eria-core/model/property"
)

type Model struct {
	Version    string
	Properties map[string]ModelProperty
	Events     map[string]ModelEvent
	Actions    map[string]ModelAction
}

type ModelProperty struct {
	DefaultValue interface{}
	Meta         propertyModel.Meta
}

type ModelEvent struct {
	Meta eventModel.Meta
}

type ModelAction struct {
	Meta actionModel.Meta
}
