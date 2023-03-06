package eria

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
	zlog "github.com/rs/zerolog/log"
)

type EriaServer struct {
	host        string
	port        uint
	exposedAddr string
	wait        sync.WaitGroup
	cancel      context.CancelFunc
	ctx         context.Context
	things      map[string]*EriaThing
	*producer.Producer
}

func NewServer(host string, port uint, exposedAddr string, instance string) *EriaServer {
	ctx, cancel := context.WithCancel(context.Background())

	server := &EriaServer{
		host:        host,
		port:        port,
		exposedAddr: exposedAddr,
		ctx:         ctx,
		cancel:      cancel,
		things:      map[string]*EriaThing{},
	}
	p := producer.New(&server.wait)
	server.Producer = p

	server.AddAppController(instance)

	return server
}

func NewThingDescription(urn string, tdVersion string, title string, description string, capabilities []string) (*thing.Thing, error) {
	td, err := NewThingFromModels(
		urn,
		tdVersion,
		title,
		description,
		capabilities,
	)

	if err != nil {
		return nil, err
	}
	//Add Versions
	td.AddContext("schema", "https://schema.org/")
	td.AddVersion("schema:softwareVersion", AppVersion)

	// Add Security
	noSecurityScheme := securityScheme.NewNoSecurity()
	td.AddSecurity("no_sec", noSecurityScheme)

	return td, err
}

func (s *EriaServer) AddThing(ref string, td *thing.Thing) (*EriaThing, error) {
	exposedThing := s.Produce(ref, td)
	eriaThing := &EriaThing{
		ref:              ref,
		propertyHandlers: map[string]PropertyData{},
		ExposedThing:     exposedThing,
	}
	for key, property := range td.Properties {
		property := property // Copy https://go.dev/doc/faq#closures_and_goroutines
		var propertyData PropertyData
		switch property.Type {
		case "boolean":
			propertyData = &PropertyBooleanData{
				value: property.Data.Default.(bool),
			}
		case "integer":
			propertyData = &PropertyIntegerData{
				value: property.Data.Default.(int),
			}
		case "number":
			propertyData = &PropertyNumberData{
				value: property.Data.Default.(float64),
			}
		case "string":
			propertyData = &PropertyStringData{
				value: property.Data.Default.(string),
			}
		}

		eriaThing.propertyHandlers[key] = propertyData
		exposedThing.SetPropertyReadHandler(key, getPropertyReadHandler(propertyData))
		exposedThing.SetPropertyWriteHandler(key, getPropertyWriteHandler(propertyData))
	}
	s.things[ref] = eriaThing
	return eriaThing, nil
}

func getPropertyReadHandler(propertyData PropertyData) producer.PropertyReadHandler {
	return func(t *producer.ExposedThing, name string) (interface{}, error) {
		if _, ok := t.ExposedProperties[name]; ok {
			value := propertyData.Get()
			zlog.Trace().Str("property", name).Interface("value", value).Msg("[core:propertyReadHandler] Value get")
			return value, nil
		}
		return nil, fmt.Errorf("property %s not found", name)
	}
}

func getPropertyWriteHandler(propertyData PropertyData) producer.PropertyWriteHandler {
	return func(t *producer.ExposedThing, name string, value interface{}) error {
		if property, ok := t.ExposedProperties[name]; ok {
			if err := property.Data.Check(value); err != nil {
				zlog.Error().Str("property", name).Interface("value", value).Err(err).Msg("[core:propertyWriteHandler]")
				return err
			}
			_, err := propertyData.Set(value)
			if err != nil {
				zlog.Error().Str("property", name).Interface("value", value).Err(err).Msg("[core:propertyWriteHandler]")
				return err
			}
			zlog.Trace().Str("property", name).Interface("value", value).Msg("[core:propertyWriteHandler] Value set")
			if err := t.EmitPropertyChange(name); err != nil {
				zlog.Error().Str("property", name).Interface("value", value).Err(err).Msg("[core:propertyWriteHandler]")
				return err
			}
			return nil
		}
		return fmt.Errorf("property %s not found", name)
	}
}

// TODO Remove
// func (s *EriaServer) SetPropertyValue(ref string, property string, value interface{}) bool {
// 	if thing, in := s.things[ref]; in {
// 		return thing.SetPropertyValue(property, value)
// 	} else {
// 		zlog.Error().Str("thing", ref).Msg("[core] thing not found")
// 	}
// 	return false
// }

func (s *EriaServer) StartServer() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	httpServer := protocolHttp.NewServer(addr, s.exposedAddr, _appName, _appName+" "+AppVersion)
	s.AddServer(httpServer)
	wsServer := protocolWebSocket.NewServer(httpServer)
	s.AddServer(wsServer)

	s.Expose()
	go func() {
		<-s.ctx.Done()
		s.Stop()
		//		_wait.Done()
	}()

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// Block until keyboard interrupt is received.
	<-c
	zlog.Info().Msg("[core:WaitForSignal] Keyboard interrupt received, Stopping...")
	s.cancel()
	// Wait for the child goroutine to finish, which will only occur when
	// the child process has stopped and the call to cmd.Wait has returned.
	// This prevents main() exiting prematurely.
	s.wait.Wait()
}

// SetThing register a thing
func SetThing(t *thing.Thing) (*producer.Producer, *producer.ExposedThing) {
	//	_wait.Add(1)
	p := producer.New(nil)
	exposedThing := p.Produce("", t)

	return p, exposedThing
}

// SetThings register a list of things and launch servers
func SetThings(ts map[string]*thing.Thing) (*producer.Producer, map[string]*producer.ExposedThing) {
	exposedThings := make(map[string]*producer.ExposedThing)
	//	_wait.Add(1)
	p := producer.New(nil)
	for ref, t := range ts {
		exposedThings[ref] = p.Produce(ref, t)
	}
	return p, exposedThings
}
