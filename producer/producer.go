package eriaproducer

import (
	"context"
	"fmt"
	"sync"

	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
	"github.com/project-eria/go-wot/thing"
)

type EriaProducer struct {
	host                        string
	port                        uint
	exposedAddr                 string
	appName                     string
	appVersion                  string
	buildDate                   string
	wait                        sync.WaitGroup
	cancel                      context.CancelFunc
	ctx                         context.Context
	things                      map[string]producer.ExposedThing
	propertyDefaultDataHandlers map[string]map[string]*PropertyData
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

func (p *EriaProducer) StartServer() {
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	httpServer := protocolHttp.NewServer(addr, p.exposedAddr, p.appName, fmt.Sprintf("%s %s (%s)", p.appName, p.appVersion, p.buildDate))
	wsServer := protocolWebSocket.NewServer(httpServer)
	// wsServer Needs to be added BEFORE httpServer,
	// in order to call the .Use(WS) middleware, before the .Get(HTTP)
	p.AddServer(wsServer)
	p.AddServer(httpServer)

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
