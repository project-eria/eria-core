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
		propertyHandlers: map[string]*PropertyData{},
		ExposedThing:     exposedThing,
	}
	for key, property := range td.Properties {
		property := property // Copy https://go.dev/doc/faq#closures_and_goroutines
		var propertyData = &PropertyData{
			value:     property.Data.Default,
			valueType: property.Type,
		}

		eriaThing.propertyHandlers[key] = propertyData
		exposedThing.SetPropertyReadHandler(key, getPropertyDefaultReadHandler(propertyData))
		exposedThing.SetPropertyWriteHandler(key, getPropertyDefaultWriteHandler(propertyData))
	}
	s.things[ref] = eriaThing
	return eriaThing, nil
}

func getPropertyDefaultReadHandler(propertyData *PropertyData) producer.PropertyReadHandler {
	return func(t *producer.ExposedThing, name string, params map[string]string) (interface{}, error) {
		value := propertyData.Get()
		zlog.Trace().Str("property", name).Interface("value", value).Msg("[core:propertyReadHandler] Value get")
		return value, nil
	}
}

func getPropertyDefaultWriteHandler(propertyData *PropertyData) producer.PropertyWriteHandler {
	return func(t *producer.ExposedThing, name string, value interface{}, params map[string]string) error {
		_, err := propertyData.Set(value)
		if err != nil {
			zlog.Error().Str("property", name).Interface("value", value).Err(err).Msg("[core:propertyWriteHandler]")
			return err
		}
		zlog.Trace().Str("property", name).Interface("value", value).Msg("[core:propertyWriteHandler] Value set")
		return nil
	}
}

func (s *EriaServer) StartServer() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	httpServer := protocolHttp.NewServer(addr, s.exposedAddr, _appName, fmt.Sprintf("%s %s (%s)", _appName, AppVersion, BuildDate))
	wsServer := protocolWebSocket.NewServer(httpServer)
	// wsServer Needs to be added BEFORE httpServer,
	// in order to call the .Use(WS) middleware, before the .Get(HTTP)
	s.AddServer(wsServer)
	s.AddServer(httpServer)

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
