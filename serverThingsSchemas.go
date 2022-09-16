package eria

import (
	"errors"

	"github.com/project-eria/eria-core/model"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// NewFromSchemas return a thing from schemas @type
func NewThingFromSchemas(urn string, version string, title string, description string, capabilities []string) (*thing.Thing, error) {
	t, err := thing.New(urn, version, title, description, nil)
	if err != nil {
		return nil, err
	}

	if err := AddSchemas(t, capabilities, ""); err != nil {
		return nil, err
	}

	return t, nil
}

// AddSchemas add capabilities to a thing using schemas @type
func AddSchemas(t *thing.Thing, capabilities []string, postfix string) error {
	// TODO (remove?) t.SetContext("http://www.w3.org/ns/td")

	for _, capability := range capabilities {
		if err := AddSchema(t, capability, postfix); err != nil {
			return err
		}
	}

	return nil
}

// AddSchema add a capability to a thing using schema @type
func AddSchema(t *thing.Thing, capability string, postfix string) error {
	log.Info().Str("capability", capability).Msg("[thing:AddSchema] Adding capability")
	propertiesSchemas, inProperties := model.CapabilitiesProperties[capability]
	if inProperties {
		for _, schema := range propertiesSchemas {
			id := schema.Id + postfix
			if _, err := AddPropertyFromSchema(t, id, schema.DefaultValue, schema.Name); err != nil {
				return err
			}
		}
	}
	actionsSchemas, inActions := model.CapabilitiesActions[capability]
	if inActions {
		for _, schema := range actionsSchemas {
			id := schema.Id + postfix
			if _, err := AddActionFromSchema(t, id, schema.Name); err != nil {
				return err
			}
		}
	}
	if !inProperties && !inActions {
		return errors.New("Capability schema '" + capability + "' not found")
	}

	// TODO (remove?) t.AddType(capability)

	return nil
}
