package eriaconsumer

import (
	"errors"

	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

type Property interface {
	Value() (interface{}, error)
	Observe(listener PropertyObserver) (uint16, error)
	UnObserve(ref uint16) error
}

type EriaProperty struct {
	key          string
	uriVariables map[string]interface{}
	onError      func(*EriaProperty, error)
	thing        consumer.ConsumedThing
	observers    []PropertyObserver
}

type PropertyOption func(*EriaProperty)
type PropertyObserver func(value interface{}, err error)

func PropertyOnError(onError func(*EriaProperty, error)) PropertyOption {
	return func(t *EriaProperty) {
		t.onError = onError
	}
}

func PropertyUriVariables(uriVariables map[string]interface{}) PropertyOption {
	return func(t *EriaProperty) {
		t.uriVariables = uriVariables
	}
}

func (t *EriaThing) Property(key string, opts ...PropertyOption) Property {
	if t == nil {
		zlog.Error().Msg("[core:consumer:Property] nil Thing")
		return nil
	}

	// TODO test if property exists ? or defer to Value/Observe?

	p := &EriaProperty{
		key:   key,
		thing: t,
	}
	// Apply options
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *EriaProperty) Value() (interface{}, error) {
	return p.thing.ReadProperty(p.key, p.uriVariables)
}

func (p *EriaProperty) Observe(listener PropertyObserver) (uint16, error) {
	var err error
	if len(p.observers) == 0 {
		// No observateur yet, we connect the property
		err = p.thing.ObserveProperty(p.key, p.uriVariables, p.observer())
	}

	if err != nil {
		return 0, err
	}
	p.observers = append(p.observers, listener)
	return uint16(len(p.observers) - 1), nil
}

func (p *EriaProperty) UnObserve(ref uint16) error {
	if p.observers[ref] == nil {
		return errors.New("invalid observer reference")
	}
	// TODO
	// - Call UnObserveProperty if len(p.observers) == 0

	// remove listener from observers
	p.observers = append(p.observers[:ref], p.observers[ref+1:]...)
	return nil
}

func (p *EriaProperty) observer() func(value interface{}, err error) {
	// We send the notification to all observers
	return func(value interface{}, err error) {
		zlog.Trace().Int("observers", len(p.observers)).Msg("[core:consumer:Property:observer] notifying observers")
		for _, observer := range p.observers {
			observer(value, err)
		}
	}
}
