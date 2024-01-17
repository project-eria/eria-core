package model

import (
	"errors"

	actionModel "github.com/project-eria/eria-core/model/action"
	eventModel "github.com/project-eria/eria-core/model/event"
	propertyModel "github.com/project-eria/eria-core/model/property"
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/thing"
	zlog "github.com/rs/zerolog/log"
)

// NewThingFromModels return a thing from schemas @type
func NewThingFromModels(urn string, version string, title string, description string, modelIds []string) (*thing.Thing, error) {
	t, err := thing.New(urn, version, title, description, nil)
	if err != nil {
		return nil, err
	}

	if err := AddModels(t, modelIds, ""); err != nil {
		return nil, err
	}

	return t, nil
}

// AddModels add capabilities to a thing using schemas @type
func AddModels(t *thing.Thing, modelIds []string, postfix string) error {
	for _, modelId := range modelIds {
		if err := AddModel(t, modelId, postfix); err != nil {
			return err
		}
	}

	return nil
}

// AddModel add a capability to a thing using schema @type
func AddModel(t *thing.Thing, modelId string, postfix string) error {
	zlog.Info().Str("model", modelId).Msg("[thing:AddModel] Adding model")
	modelType, modelExists := Models[modelId]
	if !modelExists {
		return errors.New("Model '" + modelId + "' not found")
	}

	for key, property := range modelType.Properties {
		id := key + postfix
		if _, err := AddProperty(t, id, property.Meta, property.UriVariables); err != nil {
			return err
		}
	}
	for key, event := range modelType.Events {
		id := key + postfix
		if _, err := AddEvent(t, id, event.Meta, event.UriVariables); err != nil {
			return err
		}
	}
	for key, action := range modelType.Actions {
		id := key + postfix
		if _, err := AddAction(t, id, action.Meta, action.UriVariables); err != nil {
			return err
		}
	}

	return nil
}

func AddAction(t *thing.Thing, id string, meta actionModel.Meta, uriVariables map[string]dataSchema.Data) (*interaction.Action, error) {
	zlog.Trace().Str("action", id).Msg("[thing:AddAction] Adding action")
	action := interaction.NewAction(
		id,
		meta.Title,
		meta.Description,
		meta.Options...,
	)
	// TODO (remove?) action.ATtype = schema
	t.AddAction(action)
	return action, nil
}

// AddProperty return an property from schema @type
func AddProperty(t *thing.Thing, id string, meta propertyModel.Meta, uriVariables map[string]dataSchema.Data) (*interaction.Property, error) {
	zlog.Trace().Str("property", id).Msg("[thing:AddProperty] Adding property")
	property := interaction.NewProperty(
		id,
		meta.Title,
		meta.Description,
		meta.Data,
		meta.Options...,
	)
	t.AddProperty(property)
	// TODO (remove?) property.ATtype = schema
	return property, nil
}

func AddEvent(t *thing.Thing, id string, meta eventModel.Meta, uriVariables map[string]dataSchema.Data) (*interaction.Event, error) {
	zlog.Trace().Str("event", id).Msg("[thing:AddEvent] Adding event")
	event := interaction.NewEvent(
		id,
		meta.Title,
		meta.Description,
		meta.Options...,
	)
	t.AddEvent(event)
	return event, nil
}
