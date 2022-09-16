package eria

import (
	"errors"
)

type PropertyData interface {
	Set(interface{}) (bool, error)
	Get() interface{}
}

// Boolean
type PropertyBooleanData struct {
	value bool
}

func (p *PropertyBooleanData) Set(value interface{}) (bool, error) {
	var changed bool
	if newValue, ok := value.(bool); ok {
		if p.value != newValue {
			p.value = newValue
			changed = true
		}
		return changed, nil
	} else {
		return false, errors.New("provided data is not type bool")
	}
}

func (p *PropertyBooleanData) Get() interface{} {
	return p.value
}

// Integer
type PropertyIntegerData struct {
	value int
}

func (p *PropertyIntegerData) Set(value interface{}) (bool, error) {
	var changed bool
	if newValue, ok := value.(int); ok {
		if p.value != newValue {
			p.value = newValue
			changed = true
		}
		return changed, nil
	} else {
		return false, errors.New("provided data is not type int")
	}
}

func (p *PropertyIntegerData) Get() interface{} {
	return p.value
}

// Number
type PropertyNumberData struct {
	value float64
}

func (p *PropertyNumberData) Set(value interface{}) (bool, error) {
	var changed bool
	if newValue, ok := value.(float64); ok {
		if p.value != newValue {
			p.value = newValue
			changed = true
		}
		return changed, nil
	} else {
		return false, errors.New("provided data is not type float64")
	}
}

func (p *PropertyNumberData) Get() interface{} {
	return p.value
}

// String
type PropertyStringData struct {
	value string
}

func (p *PropertyStringData) Set(value interface{}) (bool, error) {
	var changed bool
	if newValue, ok := value.(string); ok {
		if p.value != newValue {
			p.value = newValue
			changed = true
		}
		return changed, nil
	} else {
		return false, errors.New("provided data is not type string")
	}
}

func (p *PropertyStringData) Get() interface{} {
	return p.value
}
