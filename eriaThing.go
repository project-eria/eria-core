package eria

import (
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

type EriaThing struct {
	ref              string
	propertyHandlers map[string]PropertyData
	*producer.ExposedThing
}

func (t *EriaThing) SetPropertyValue(property string, value interface{}) bool {
	if t == nil {
		zlog.Error().Msg("[core:SetPropertyValue] nil thing")
		return false
	}
	if propertyData, in := t.propertyHandlers[property]; in {
		changed, err := propertyData.Set(value)
		if err != nil {
			zlog.Error().Str("thing", t.ref).Str("property", property).Err(err).Msg("[core:SetPropertyValue]")
		}
		if changed {
			zlog.Trace().Str("thing", t.ref).Str("property", property).Interface("value", value).Msg("[core:SetPropertyValue] value changed")
			t.EmitPropertyChange(property)
			return true
		}
	} else {
		zlog.Error().Str("thing", t.ref).Str("property", property).Msg("[core:SetPropertyValue] thing property handler not found")
	}
	return false
}

func (t *EriaThing) GetPropertyValue(property string) interface{} {
	if t == nil {
		zlog.Error().Msg("[core:SetPropertyValue] nil thing")
		return nil
	}

	if propertyData, in := t.propertyHandlers[property]; in {
		return propertyData.Get()
	} else {
		zlog.Error().Str("thing", t.ref).Str("property", property).Msg("[core:GetPropertyValue] thing property handler not found")
	}
	return nil
}
