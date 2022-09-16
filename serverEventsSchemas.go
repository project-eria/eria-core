package eria

import (
	"errors"

	"github.com/project-eria/eria-core/model"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

func AddEvent(t *thing.Thing, id string, meta model.EventDesc) (*interaction.Event, error) {
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

// AddEventFromSchema return an event from schema @type
func AddEventFromSchema(t *thing.Thing, id string, schema string) (*interaction.Event, error) {
	log.Info().Str("schema", schema).Msg("[thing:AddEventFromSchema] Adding event")
	if meta, in := model.Events[schema]; in {
		return AddEvent(t, id, meta)
	}
	return nil, errors.New("Event schema '" + schema + "' not found")
}
