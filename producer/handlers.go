package eriaproducer

import (
	"errors"

	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

/**
 * Set default read/write handlers for a property
 * @param {producer.ExposedThing} t - exposed thing
 * @param {string} ref - property name
 * @returns {PropertyData}, {Error}
 */
func (p *EriaProducer) PropertyUseDefaultHandlers(t producer.ExposedThing, ref string) (*PropertyData, error) {
	if property, ok := t.TD().Properties[ref]; ok {
		var propertyData = &PropertyData{
			value:     property.Data.Default,
			valueType: property.Type,
		}
		t.SetPropertyReadHandler(ref, getPropertyDefaultReadHandler(propertyData))
		t.SetPropertyWriteHandler(ref, getPropertyDefaultWriteHandler(propertyData))

		p.propertyDefaultDataHandlers[t.Ref()][ref] = propertyData

		return propertyData, nil
	}
	return nil, errors.New("property not found")
}

func getPropertyDefaultReadHandler(propertyData *PropertyData) producer.PropertyReadHandler {
	return func(t producer.ExposedThing, property string, parameters map[string]interface{}) (interface{}, error) {
		value := propertyData.Get()
		zlog.Trace().Str("property", property).Interface("value", value).Msg("[eriaproducer:propertyReadHandler] Value get")
		return value, nil
	}
}

func getPropertyDefaultWriteHandler(propertyData *PropertyData) producer.PropertyWriteHandler {
	return func(t producer.ExposedThing, property string, value interface{}, parameters map[string]interface{}) error {
		_, err := propertyData.Set(value)
		if err != nil {
			zlog.Error().Str("property", property).Interface("value", value).Err(err).Msg("[eriaproducer:propertyWriteHandler]")
			return err
		}
		zlog.Trace().Str("property", property).Interface("value", value).Msg("[eriaproducer:propertyWriteHandler] Value set")
		return nil
	}
}

func (p *EriaProducer) SetPropertyValue(t producer.ExposedThing, ref string, value interface{}) bool {
	if t == nil {
		zlog.Error().Msg("[eriaproducer:SetPropertyValue] nil thing")
		return false
	}
	if propertyData, in := p.propertyDefaultDataHandlers[t.Ref()][ref]; in {
		changed, err := propertyData.Set(value)
		if err != nil {
			zlog.Error().Str("thing", t.Ref()).Str("property", ref).Interface("value", value).Err(err).Msg("[eriaproducer:SetPropertyValue]")
		}
		if changed {
			zlog.Trace().Str("thing", t.Ref()).Str("property", ref).Interface("value", value).Msg("[eriaproducer:SetPropertyValue] value changed")
			t.EmitPropertyChange(ref, value, nil)
			return true
		}
	} else {
		zlog.Error().Str("thing", t.Ref()).Str("property", ref).Msg("[eriaproducer:SetPropertyValue] thing property handler not found")
	}
	return false
}

func (p *EriaProducer) GetPropertyValue(t producer.ExposedThing, ref string) interface{} {
	if t == nil {
		zlog.Error().Msg("[eriaproducer:GetPropertyValue] nil thing")
		return nil
	}

	if propertyData, in := p.propertyDefaultDataHandlers[t.Ref()][ref]; in {
		return propertyData.Get()
	} else {
		zlog.Error().Str("thing", t.Ref()).Str("property", ref).Msg("[eriaproducer:GetPropertyValue] thing property handler not found")
	}
	return nil
}

func (p *EriaProducer) AddChangeCallBack(t producer.ExposedThing, ref string, f func(interface{})) {
	if t == nil {
		zlog.Error().Msg("[eriaproducer:AddChangeCallBack] nil thing")
		return
	}

	if propertyData, in := p.propertyDefaultDataHandlers[t.Ref()][ref]; in {
		propertyData.AddChangeCallBack(f)
	} else {
		zlog.Error().Str("thing", t.Ref()).Str("property", ref).Msg("[eriaproducer:GetPropertyValue] thing property handler not found")
	}
}
