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
			log.Error().Str("thing", t.ref).Str("property", property).Err(err).Msg("[core:SetPropertyValue]")
		}
		if changed {
			t.EmitPropertyChange(property)
			log.Trace().Str("thing", t.ref).Str("property", property).Interface("value", value).Msg("[core:SetPropertyValue] value changed")
			return true
		}
	} else {
		log.Error().Str("thing", t.ref).Str("property", property).Msg("[core:SetPropertyValue] thing property handler not found")
	}
	return false
}

func (t *EriaThing) GetPropertyValue(property string) interface{} {
	if propertyData, in2 := t.propertyHandlers[property]; in2 {
		return propertyData.Get()
	} else {
		log.Error().Str("thing", t.ref).Str("property", property).Msg("[core:GetPropertyValue] thing property handler not found")
	}
	return nil
}
