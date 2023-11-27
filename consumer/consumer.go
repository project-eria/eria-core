package eriaconsumer

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
	"github.com/project-eria/go-wot/thing"
	zlog "github.com/rs/zerolog/log"
)

type EriaConsumer struct {
	things map[string]consumer.ConsumedThing // TODO store connected things by urls
	*consumer.Consumer
}

func New() *EriaConsumer {
	c := consumer.New()
	httpClient := protocolHttp.NewClient()
	c.AddClient(httpClient)
	wsClient := protocolWebSocket.NewClient()
	c.AddClient(wsClient)

	return &EriaConsumer{
		things:   map[string]consumer.ConsumedThing{},
		Consumer: c,
	}
}

// ConnectThing connect a remote thing WS server
func (c *EriaConsumer) ConnectThing(url string, onConnected func(consumer.ConsumedThing), onError func(error)) {
	// TODO test if already connected
	go func() {
		thing, err := c.ConnectThingBackoff(url)
		if err == nil {
			if onConnected != nil {
				onConnected(thing)
			}
		} else {
			if onError != nil {
				onError(err)
			}
		}
	}()
}

func (c *EriaConsumer) ConnectThingBackoff(url string) (consumer.ConsumedThing, error) {
	// A backoff schedule for when and how often to retry failed HTTP
	// requests. The first element is the time to wait after the
	// first failure, the second the time to wait after the second
	// failure, etc. After reaching the last element, retries stop
	// and the request is considered failed.
	var backoffSchedule = []time.Duration{
		1 * time.Second,
		5 * time.Second,
		5 * time.Second,
	}
	var err error
	var resp *http.Response
	for _, backoff := range backoffSchedule {
		resp, err = http.Get(url)
		if err == nil {
			break
		}

		zlog.Error().Str("url", url).Err(err).Msgf("[core:ConnectThing] Retrying in %v\n", backoff)
		time.Sleep(backoff)
	}

	if err != nil {
		return nil, errors.New("can't open thing url")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		zlog.Error().Str("status", resp.Status).Str("url", url).Msg("[core:ConnectThing] incorrect response")
		return nil, errors.New("incorrect response")
	}

	var td thing.Thing
	if err := json.NewDecoder(resp.Body).Decode(&td); err != nil {
		zlog.Error().Str("url", url).Err(err).Msg("[core:ConnectThing]")
		return nil, errors.New("can't decode json")
	}

	consumedThing := c.Consume(&td)
	c.things[td.ID] = consumedThing

	return consumedThing, nil
}
