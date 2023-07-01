package model

import (
	actionModel "github.com/project-eria/eria-core/model/action"
	eventModel "github.com/project-eria/eria-core/model/event"
	propertyModel "github.com/project-eria/eria-core/model/property"
	"github.com/project-eria/go-wot/dataSchema"
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
	UriVariables map[string]dataSchema.Data
}

type ModelEvent struct {
	Meta         eventModel.Meta
	UriVariables map[string]dataSchema.Data
}

type ModelAction struct {
	Meta         actionModel.Meta
	UriVariables map[string]dataSchema.Data
}
