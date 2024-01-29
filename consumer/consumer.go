package eriaconsumer

import (
	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
)

type Consumer interface {
	Thing(url string) Thing
	ThingFromTag(tag string) Thing
}

type EriaConsumer struct {
	things map[string]Thing
	*consumer.Consumer
}

func New() *EriaConsumer {
	c := consumer.New()
	httpClient := protocolHttp.NewClient()
	c.AddClient(httpClient)
	wsClient := protocolWebSocket.NewClient()
	c.AddClient(wsClient)

	return &EriaConsumer{
		things:   map[string]Thing{},
		Consumer: c,
	}
}
