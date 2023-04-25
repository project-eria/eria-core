package eria

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"

	zlog "github.com/rs/zerolog/log"
)

type PropertyData interface {
	Set(interface{}) (bool, error)
	Get() interface{}
	AddChangeCallBack(func(interface{}))
}

type PropertyGeneralData struct {
	mu              sync.RWMutex
	changeCallbacks []func(interface{})
}

func (p *PropertyGeneralData) AddChangeCallBack(f func(interface{})) {
	if p == nil {
		zlog.Error().Msg("[core:AddChangeCallBack] nil property data")
		return
	}
	p.changeCallbacks = append(p.changeCallbacks, f)
}

func (p *PropertyGeneralData) emitChangeCallback(value interface{}) {
	for _, f := range p.changeCallbacks {
		go f(value)
	}
}

// Boolean
type PropertyBooleanData struct {
	value bool
	*PropertyGeneralData
}

func (p *PropertyBooleanData) Set(value interface{}) (bool, error) {
	var changed bool
	if newValue, ok := value.(bool); ok {
		if p.value != newValue {
			p.value = newValue
			changed = true
			p.emitChangeCallback(newValue)
		}
		return changed, nil
	} else {
		return false, errors.New("provided data is not type bool")
	}
}

func (p *PropertyBooleanData) Get() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.value
}

// Integer
type PropertyIntegerData struct {
	value int
	*PropertyGeneralData
}

func (p *PropertyIntegerData) Set(value interface{}) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var changed bool
	if newValue, ok := value.(int); ok {
		if p.value != newValue {
			p.value = newValue
			changed = true
			p.emitChangeCallback(newValue)
		}
		return changed, nil
	} else {
		return false, errors.New("provided data is not type int")
	}
}

func (p *PropertyIntegerData) Get() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.value
}

// Number
type PropertyNumberData struct {
	value float64
	*PropertyGeneralData
}

func (p *PropertyNumberData) Set(value interface{}) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var changed bool
	if newValue, ok := value.(float64); ok {
		if p.value != newValue {
			p.value = newValue
			changed = true
			p.emitChangeCallback(newValue)
		}
		return changed, nil
	} else {
		return false, errors.New("provided data is not type float64")
	}
}

func (p *PropertyNumberData) Get() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.value
}

// String
type PropertyStringData struct {
	value string
	*PropertyGeneralData
}

func (p *PropertyStringData) Set(value interface{}) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var changed bool
	if newValue, ok := value.(string); ok {
		if p.value != newValue {
			p.value = newValue
			changed = true
			p.emitChangeCallback(newValue)
		}
		return changed, nil
	} else {
		return false, errors.New("provided data is not type string")
	}
}

func (p *PropertyStringData) Get() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.value
}

// Object
type PropertyObjectData struct {
	value map[string]interface{}
	*PropertyGeneralData
}

func (p *PropertyObjectData) Set(value interface{}) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var changed bool
	if newValue, ok := value.(map[string]interface{}); ok {
		newJson, err := json.Marshal(newValue)
		if err != nil {
			return false, errors.New("provided data can't be decoded")
		}
		oldJson, _ := json.Marshal(p.value)
		if !bytes.Equal(newJson, oldJson) {
			p.value = newValue
			changed = true
			p.emitChangeCallback(newValue)
		}
		return changed, nil
	} else {
		return false, errors.New("provided data is not type object")
	}
}

func (p *PropertyObjectData) Get() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.value
}
