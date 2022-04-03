package eria

import (
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

type EriaThing struct {
	ref              string
	propertyHandlers map[string]PropertyData
	*producer.ExposedThing
}

func (t *EriaThing) SetPropertyValue(property string, value interface{}) bool {
	if propertyData, in2 := t.propertyHandlers[property]; in2 {
		changed, err := propertyData.Set(value)
		if err != nil {
			log.Error().Str("thing", t.ref).Str("property", property).Err(err).Msg("[core]")
		}
		if changed {
			t.EmitPropertyChange(property)
			log.Trace().Str("thing", t.ref).Str("property", property).Interface("value", value).Msg("[core] value changed")
			return true
		}
	} else {
		log.Error().Str("thing", t.ref).Str("property", property).Msg("[core] thing property handler not found")
	}
	return false
}
