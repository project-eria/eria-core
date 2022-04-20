package eria

import (
	"errors"

	"github.com/project-eria/eria-core/model"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

func AddAction(t *thing.Thing, id string, meta model.ActionDesc) (*interaction.Action, error) {
	log.Trace().Msg("[thing:AddAction] Adding action")
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

// AddActionFromSchema return an action from schema @type
func AddActionFromSchema(t *thing.Thing, id string, schema string) (*interaction.Action, error) {
	log.Info().Str("schema", schema).Msg("[thing:AddActionFromSchema] Adding action")
	if meta, in := model.Actions[schema]; in {
		return AddAction(t, id, meta)
	}
	return nil, errors.New("Action schema '" + schema + "' not found")
}
