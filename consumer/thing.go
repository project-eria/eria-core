package eriaconsumer

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

type Thing interface {
	Property(key string, opts ...PropertyOption) Property
}

type EriaThing struct {
	Tags []string
	// MdnsInstance is the DNS-SD identifier the Thing was discovered with,
	// kept for re-discovery after a connection loss (the IP may change).
	MdnsInstance   string
	url            string
	onConnected    func(*EriaThing)
	onError        func(*EriaThing, error)
	onDisconnected func(*EriaThing, error)
	method         ConnectMethod
	reconnecting   atomic.Bool
	mu             sync.RWMutex // guards the ConsumedThing swap on reconnection
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

// WithThingOnDisconnected sets the callback invoked (once per loss) when an
// interaction with the Thing fails, hinting that the connection is lost.
func WithThingOnDisconnected(onDisconnected func(*EriaThing, error)) ThingOption {
	return func(t *EriaThing) {
		t.onDisconnected = onDisconnected
	}
}

// WithThingMdnsInstance records the DNS-SD identifier used to discover the Thing.
func WithThingMdnsInstance(instance string) ThingOption {
	return func(t *EriaThing) {
		t.MdnsInstance = instance
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
		if err != nil {
			if eriaThing.onError != nil {
				eriaThing.onError(nil, err)
			}
			return
		}
		thing := c.Consume(td)
		eriaThing.ConsumedThing = thing
		eriaThing.url = url
		c.thingsMu.Lock()
		c.things[url] = eriaThing
		c.thingsMu.Unlock()
		if eriaThing.onConnected != nil {
			eriaThing.onConnected(eriaThing)
		}
	}()
}

// ReconnectThing re-fetches the TD at url and swaps the underlying
// ConsumedThing, preserving the EriaThing pointer held by the app. It retries
// until the Thing answers again. Note: WebSocket observations established
// before the loss keep reconnecting to the previous address; reads, writes
// and actions use the new one.
func (c *EriaConsumer) ReconnectThing(t *EriaThing, url string) {
	for {
		td, err := t.method(url)
		if err == nil {
			thing := c.Consume(td)
			t.mu.Lock()
			oldURL := t.url
			t.ConsumedThing = thing
			t.url = url
			t.mu.Unlock()
			if oldURL != url {
				c.thingsMu.Lock()
				delete(c.things, oldURL)
				c.things[url] = t
				c.thingsMu.Unlock()
			}
			t.reconnecting.Store(false)
			zlog.Info().Str("url", url).Msg("[core:consumer:ReconnectThing] Thing reconnected")
			if t.onConnected != nil {
				t.onConnected(t)
			}
			return
		}
		zlog.Warn().Str("url", url).Err(err).Msg("[core:consumer:ReconnectThing] still unreachable, retrying in 30s")
		time.Sleep(30 * time.Second)
	}
}

// notifyLoss reports a failed interaction. The onDisconnected callback is
// fired once: further losses are ignored until a reconnection resets the flag.
func (t *EriaThing) notifyLoss(err error) {
	if t.onDisconnected == nil {
		return
	}
	if !t.reconnecting.CompareAndSwap(false, true) {
		return // a loss is already being handled
	}
	go t.onDisconnected(t, err)
}

// The interaction methods below shadow the embedded ConsumedThing to
// (a) read it under lock (it is swapped on reconnection) and (b) report
// interaction errors as a possible connection loss.

func (t *EriaThing) consumedThing() consumer.ConsumedThing {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.ConsumedThing
}

func (t *EriaThing) ReadProperty(name string, uriVariables map[string]interface{}) (interface{}, error) {
	value, err := t.consumedThing().ReadProperty(name, uriVariables)
	if err != nil {
		t.notifyLoss(err)
	}
	return value, err
}

func (t *EriaThing) WriteProperty(name string, uriVariables map[string]interface{}, value interface{}) (interface{}, error) {
	result, err := t.consumedThing().WriteProperty(name, uriVariables, value)
	if err != nil {
		t.notifyLoss(err)
	}
	return result, err
}

func (t *EriaThing) InvokeAction(name string, uriVariables map[string]interface{}, value interface{}) (interface{}, error) {
	result, err := t.consumedThing().InvokeAction(name, uriVariables, value)
	if err != nil {
		t.notifyLoss(err)
	}
	return result, err
}

func (c *EriaConsumer) Thing(url string) Thing {
	c.thingsMu.RLock()
	defer c.thingsMu.RUnlock()
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
	c.thingsMu.RLock()
	defer c.thingsMu.RUnlock()
	for _, t := range c.things {
		et := t.(*EriaThing)
		if arrutil.StringsHas(et.Tags, tag) {
			return t
		}
	}
	zlog.Error().Msg("[core:consumer:ThingFromTag] Can't find Thing")
	return nil
}
