package eria

import (
	"errors"

	"github.com/project-eria/eria-core/model"
	actionModel "github.com/project-eria/eria-core/model/action"
	eventModel "github.com/project-eria/eria-core/model/event"
	propertyModel "github.com/project-eria/eria-core/model/property"
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
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
	log.Info().Str("model", modelId).Msg("[thing:AddModel] Adding model")
	modelType, modelExists := model.Models[modelId]
	if !modelExists {
		return errors.New("Model '" + modelId + "' not found")
	}

	for key, property := range modelType.Properties {
		id := key + postfix
		if _, err := AddProperty(t, id, property.DefaultValue, property.Meta); err != nil {
			return err
		}
	}
	for key, event := range modelType.Events {
		id := key + postfix
		if _, err := AddEvent(t, id, event.Meta); err != nil {
			return err
		}
	}
	for key, action := range modelType.Actions {
		id := key + postfix
		if _, err := AddAction(t, id, action.Meta); err != nil {
			return err
		}
	}

	return nil
}

func AddAction(t *thing.Thing, id string, meta actionModel.Meta) (*interaction.Action, error) {
	log.Trace().Str("action", id).Msg("[thing:AddAction] Adding action")
	action := interaction.NewAction(
		id,
		meta.Title,
		meta.Description,
		nil,
		nil,
	)
	// TODO (remove?) action.ATtype = schema
	t.AddAction(action)
	return action, nil
}

// AddProperty return an property from schema @type
func AddProperty(t *thing.Thing, id string, defaultValue interface{}, meta propertyModel.Meta) (*interaction.Property, error) {
	log.Trace().Str("property", id).Msg("[thing:AddProperty] Adding property")
	var data dataSchema.Data
	switch meta.Type {
	case "boolean":
		data = dataSchema.NewBoolean(defaultValue.(bool))
	case "integer":
		data = dataSchema.NewInteger(defaultValue.(int), meta.Unit, meta.Minimum, meta.Maximum)
	case "number":
		data = dataSchema.NewNumber(defaultValue.(float64), meta.Unit, meta.Minimum, meta.Maximum)
	case "string":
		data = dataSchema.NewString(defaultValue.(string), meta.MinLength, meta.MaxLength, meta.Pattern)
	}
	property := interaction.NewProperty(
		id,
		meta.Title,
		meta.Description,
		false,
		false,
		true,
		data,
	)
	t.AddProperty(property)
	// TODO (remove?) property.ATtype = schema
	return property, nil
}

func AddEvent(t *thing.Thing, id string, meta eventModel.Meta) (*interaction.Event, error) {
	log.Trace().Str("event", id).Msg("[thing:AddEvent] Adding event")
	event := interaction.NewEvent(
		id,
		meta.Title,
		meta.Description,
		nil,
	)
	t.AddEvent(event)
	return event, nil
}
