package eriaproducer

import (
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

type EriaThing struct {
	ref                         string
	propertyDefaultDataHandlers map[string]*PropertyData
	*producer.ExposedThing
}

func NewEriaThing(ref string, exposedThing *producer.ExposedThing) *EriaThing {
	return &EriaThing{
		ref:                         ref,
		propertyDefaultDataHandlers: map[string]*PropertyData{},
		ExposedThing:                exposedThing,
	}
}

func (t *EriaThing) PropertyUseDefaultHandlers(name string) {
	if property, ok := t.Td.Properties[name]; ok {
		var propertyData = &PropertyData{
			value:     property.Data.Default,
			valueType: property.Type,
		}
		t.propertyDefaultDataHandlers[name] = propertyData
		t.SetPropertyReadHandler(name, getPropertyDefaultReadHandler(propertyData))
		t.SetPropertyWriteHandler(name, getPropertyDefaultWriteHandler(propertyData))
	} else {
		zlog.Error().Str("property", name).Msg("[core:PropertyUseDefaultHandlers] Property not found")
	}
}

func getPropertyDefaultReadHandler(propertyData *PropertyData) producer.PropertyReadHandler {
	return func(t *producer.ExposedThing, name string, options map[string]string) (interface{}, error) {
		value := propertyData.Get()
		zlog.Trace().Str("property", name).Interface("value", value).Msg("[core:propertyReadHandler] Value get")
		return value, nil
	}
}

func getPropertyDefaultWriteHandler(propertyData *PropertyData) producer.PropertyWriteHandler {
	return func(t *producer.ExposedThing, name string, value interface{}, options map[string]string) error {
		_, err := propertyData.Set(value)
		if err != nil {
			zlog.Error().Str("property", name).Interface("value", value).Err(err).Msg("[core:propertyWriteHandler]")
			return err
		}
		zlog.Trace().Str("property", name).Interface("value", value).Msg("[core:propertyWriteHandler] Value set")
		return nil
	}
}

func (t *EriaThing) SetPropertyValue(property string, value interface{}) bool {
	if t == nil {
		zlog.Error().Msg("[core:SetPropertyValue] nil thing")
		return false
	}
	if propertyData, in := t.propertyDefaultDataHandlers[property]; in {
		changed, err := propertyData.Set(value)
		if err != nil {
			zlog.Error().Str("thing", t.ref).Str("property", property).Interface("value", value).Err(err).Msg("[core:SetPropertyValue]")
		}
		if changed {
			zlog.Trace().Str("thing", t.ref).Str("property", property).Interface("value", value).Msg("[core:SetPropertyValue] value changed")
			t.EmitPropertyChange(property, value, map[string]string{})
			return true
		}
	} else {
		zlog.Error().Str("thing", t.ref).Str("property", property).Msg("[core:SetPropertyValue] thing property handler not found")
	}
	return false
}

func (t *EriaThing) GetPropertyValue(property string) interface{} {
	if t == nil {
		zlog.Error().Msg("[core:GetPropertyValue] nil thing")
		return nil
	}

	if propertyData, in := t.propertyDefaultDataHandlers[property]; in {
		return propertyData.Get()
	} else {
		zlog.Error().Str("thing", t.ref).Str("property", property).Msg("[core:GetPropertyValue] thing property handler not found")
	}
	return nil
}

func (t *EriaThing) AddChangeCallBack(property string, f func(interface{})) {
	if t == nil {
		zlog.Error().Msg("[core:AddChangeCallBack] nil thing")
		return
	}

	if propertyData, in := t.propertyDefaultDataHandlers[property]; in {
		propertyData.AddChangeCallBack(f)
	} else {
		zlog.Error().Str("thing", t.ref).Str("property", property).Msg("[core:GetPropertyValue] thing property handler not found")
	}
}
