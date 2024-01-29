package eriaconsumer

import (
	"github.com/gookit/goutil/arrutil"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

type Thing interface {
	Property(key string, opts ...PropertyOption) Property
}

type EriaThing struct {
	Tags        []string
	onConnected func(*EriaThing)
	onError     func(*EriaThing, error)
	method      ConnectMethod
	consumer.ConsumedThing
}

type ThingOption func(*EriaThing)

func WithThingOnConnected(onConnected func(*EriaThing)) ThingOption {
	return func(t *EriaThing) {
		t.onConnected = onConnected
	}
}

func WithThingOnError(onError func(*EriaThing, error)) ThingOption {
	return func(t *EriaThing) {
		t.onError = onError
	}
}

func WithThingBackoffMethod() ThingOption {
	return func(t *EriaThing) {
		t.method = BackoffMethod
	}
}

func WithThingTags(tags []string) ThingOption {
	return func(t *EriaThing) {
		t.Tags = append(t.Tags, tags...)
	}
}

// ConnectThing connect a remote thing server
func (c *EriaConsumer) ConnectThing(url string, opts ...ThingOption) {
	eriaThing := &EriaThing{
		method: BackoffMethod, // Default method
	}

	// Apply options
	for _, opt := range opts {
		opt(eriaThing)
	}

	// TODO test if already connected
	go func() {
		td, err := eriaThing.method(url)
		thing := c.Consume(td)
		eriaThing.ConsumedThing = thing
		if err == nil {
			c.things[url] = eriaThing
			if eriaThing.onConnected != nil {
				eriaThing.onConnected(eriaThing)
			}
		} else {
			if eriaThing.onError != nil {
				eriaThing.onError(nil, err)
			}
		}
	}()
}

func (c *EriaConsumer) Thing(url string) Thing {
	if t, ok := c.things[url]; ok {
		return t
	}
	zlog.Error().Msg("[core:consumer:Thing] Can't find Thing")
	return nil
}

/**
 * Returns the first Thing with the given tag
 */
func (c *EriaConsumer) ThingFromTag(tag string) Thing {
	for _, t := range c.things {
		et := t.(*EriaThing)
		if arrutil.StringsHas(et.Tags, tag) {
			return t
		}
	}
	zlog.Error().Msg("[core:consumer:ThingFromTag] Can't find Thing")
	return nil
}
