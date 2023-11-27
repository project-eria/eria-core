package eriaproducer

import (
	"sync"

	zlog "github.com/rs/zerolog/log"
)

type PropertyData struct {
	value           interface{}
	valueType       string
	mu              sync.RWMutex
	changeCallbacks []func(interface{})
}

func (p *PropertyData) AddChangeCallBack(f func(interface{})) {
	if p == nil {
		zlog.Error().Msg("[eriaproducer:AddChangeCallBack] nil property data")
		return
	}
	p.changeCallbacks = append(p.changeCallbacks, f)
}

func (p *PropertyData) emitChangeCallback(value interface{}) {
	for _, f := range p.changeCallbacks {
		go f(value)
	}
}

func (p *PropertyData) Set(value interface{}) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// TODO check if value matches data criterias
	var changed = p.valueType == "object" || p.valueType == "array" || p.value != value // We don't compare complex objects for changes
	if changed {
		p.value = value
		p.emitChangeCallback(value)
	}
	return changed, nil
}

func (p *PropertyData) Get() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.value
}

// // Object

// func (p *PropertyObjectData) Set(value interface{}) (bool, error) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	var changed bool
// 	if newValue, ok := value.(map[string]interface{}); ok {
// 		newJson, err := json.Marshal(newValue)
// 		if err != nil {
// 			return false, errors.New("provided data can't be decoded")
// 		}
// 		oldJson, _ := json.Marshal(p.value)
// 		if !bytes.Equal(newJson, oldJson) {
// 			p.value = newValue
// 			changed = true
// 			p.emitChangeCallback(newValue)
// 		}
// 		return changed, nil
// 	} else {
// 		return false, errors.New("provided data is not type object")
// 	}
// }
