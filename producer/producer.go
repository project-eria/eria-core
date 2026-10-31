package eriaproducer

import (
	"context"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolDnssd"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
	"github.com/project-eria/go-wot/thing"
	zlog "github.com/rs/zerolog/log"
)

// ControllerRef is the ref of the App Controller Thing: served under /eria,
// never announced over mDNS
const ControllerRef = "eria"

type EriaProducer struct {
	host                        string
	port                        uint
	exposedAddr                 string
	appName                     string
	appVersion                  string
	buildDate                   string
	mdnsEnabled                 bool
	mdnsScheme                  string
	mdnsNetwork                 string
	wait                        sync.WaitGroup
	cancel                      context.CancelFunc
	ctx                         context.Context
	things                      map[string]producer.ExposedThing
	propertyDefaultDataHandlers map[string]map[string]*PropertyData
	guard                       fiber.Handler
	*producer.Producer
}

func New(host string, port uint, exposedAddr string, appName string, appVersion string, buildDate string) *EriaProducer {
	ctx, cancel := context.WithCancel(context.Background())

	eriaProducer := &EriaProducer{
		host:                        host,
		port:                        port,
		exposedAddr:                 exposedAddr,
		appName:                     appName,
		appVersion:                  appVersion,
		buildDate:                   buildDate,
		ctx:                         ctx,
		cancel:                      cancel,
		things:                      map[string]producer.ExposedThing{},
		propertyDefaultDataHandlers: map[string]map[string]*PropertyData{},
	}
	p := producer.New(&eriaProducer.wait)
	eriaProducer.Producer = p

	return eriaProducer
}

// SetMdns configures the DNS-SD announcement of the exposed Things, on the
// service type of the network (`_wot._tcp` for the default one). Must be
// called before StartServer.
func (p *EriaProducer) SetMdns(enabled bool, scheme string, network string) {
	p.mdnsEnabled = enabled
	p.mdnsScheme = scheme
	p.mdnsNetwork = network
}

// Guard runs h before every HTTP and WebSocket handler of the exposed
// Things (the eria-manager master checks its admin token there). Must be
// called before StartServer.
func (p *EriaProducer) Guard(h fiber.Handler) {
	p.guard = h
}

func (p *EriaProducer) StartServer() {
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	httpServer := protocolHttp.NewServer(addr, p.exposedAddr, p.appName, fmt.Sprintf("%s %s (%s)", p.appName, p.appVersion, p.buildDate))
	if p.guard != nil {
		httpServer.Use(p.guard) // before the WebSocket upgrade and the routes
	}
	wsServer := protocolWebSocket.NewServer(httpServer)
	// wsServer Needs to be added BEFORE httpServer,
	// in order to call the .Use(WS) middleware, before the .Get(HTTP)
	p.AddServer(wsServer)
	p.AddServer(httpServer)

	if p.mdnsEnabled {
		if p.exposedAddr != "" {
			// The mDNS announcement points to the local primary IP, not the
			// exposed (proxy) address: consumers behind the proxy should use
			// the static URL config instead.
			zlog.Warn().Str("exposedAddr", p.exposedAddr).Msg("[producer:StartServer] exposedAddr is set: the mDNS announcement will point to the local address, consider mdns.disabled: true")
		}
		advertiser, err := protocolDnssd.NewAdvertiser(protocolDnssd.Options{
			Port:    int(p.port),
			Scheme:  p.mdnsScheme,
			Network: p.mdnsNetwork,
			// The App Controller (ref "eria") is not announced, like on the
			// ESP32: consumers reach it at /eria next to the app's own Thing
			// (the Home Assistant integration probes it there)
			Skip: func(ref string) bool { return ref == ControllerRef },
		})
		if err != nil {
			zlog.Error().Err(err).Msg("[producer:StartServer] can't create the mDNS advertiser, announcement disabled")
		} else {
			p.AddServer(advertiser)
		}
	}

	p.Expose()
	p.Start()
	go func() {
		<-p.ctx.Done()
		p.Stop()
		//		_wait.Done()
	}()

}

func (p *EriaProducer) StopServer() {
	p.cancel()
	// Wait for the child goroutine to finish, which will only occur when
	// the child process has stopped and the call to cmd.Wait has returned.
	// This prevents main() exiting prematurely.
	p.wait.Wait()
}

func (p *EriaProducer) AddThing(ref string, td *thing.Thing) (producer.ExposedThing, error) {
	// Add Versions
	td.AddVersion("schema:softwareVersion", p.appVersion)

	exposedThing := p.Produce(ref, td)
	p.things[ref] = exposedThing
	// Initialize the list of default properties handlers
	p.propertyDefaultDataHandlers[ref] = map[string]*PropertyData{}
	return exposedThing, nil
}

// SetThing register a thing
func SetThing(t *thing.Thing) (*producer.Producer, producer.ExposedThing) {
	//	_wait.Add(1)
	p := producer.New(nil)
	exposedThing := p.Produce("", t)

	return p, exposedThing
}

// SetThings register a list of things and launch servers
func SetThings(ts map[string]*thing.Thing) (*producer.Producer, map[string]producer.ExposedThing) {
	exposedThings := make(map[string]producer.ExposedThing)
	//	_wait.Add(1)
	p := producer.New(nil)
	for ref, t := range ts {
		exposedThings[ref] = p.Produce(ref, t)
	}
	return p, exposedThings
}

func (p *EriaProducer) GetThings() map[string]producer.ExposedThing {
	return p.things
}
