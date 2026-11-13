package eriaconsumer

import (
	"sync"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
)

type Consumer interface {
	Thing(url string) Thing
	ThingFromTag(tag string) Thing
}

type EriaConsumer struct {
	thingsMu sync.RWMutex // things is written by the ConnectThing goroutines
	things   map[string]Thing
	*consumer.Consumer
}

// Token is the app's service token sent to url, "" for none (set by
// eria-core when the app has a key, for the Things of the coordinator's inventory).
var Token = func(url string) string { return "" }

func token(url string) string { return Token(url) }

func New() *EriaConsumer {
	c := consumer.New()
	httpClient := protocolHttp.NewClient()
	httpClient.Token = token
	c.AddClient(httpClient)
	wsClient := protocolWebSocket.NewClient()
	wsClient.Token = token
	c.AddClient(wsClient)

	return &EriaConsumer{
		things:   map[string]Thing{},
		Consumer: c,
	}
}
