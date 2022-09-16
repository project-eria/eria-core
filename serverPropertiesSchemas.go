package eria

import (
	"errors"

	"github.com/project-eria/eria-core/model"
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// AddProperty return an property from schema @type
func AddProperty(t *thing.Thing, id string, defaultValue interface{}, meta model.PropertyDesc) (*interaction.Property, error) {
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

// AddPropertyFromSchema return an property from schema @type
func AddPropertyFromSchema(t *thing.Thing, id string, defaultValue interface{}, schema string) (*interaction.Property, error) {
	log.Info().Str("schema", schema).Msg("[thing:AddPropertyFromSchema] Adding property")
	if meta, in := model.Properties[schema]; in {
		return AddProperty(t, id, defaultValue, meta)
	}
	return nil, errors.New("Property schema '" + schema + "' not found")
}
