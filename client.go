package eria

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

type EriaClient struct {
	things map[string]*consumer.ConsumedThing
	*consumer.Consumer
}

func NewClient() *EriaClient {
	c := consumer.New()
	httpClient := protocolHttp.NewClient()
	c.AddClient(httpClient)
	wsClient := protocolWebSocket.NewClient()
	c.AddClient(wsClient)

	return &EriaClient{
		things:   map[string]*consumer.ConsumedThing{},
		Consumer: c,
	}
}

// ConsumeThing connect a remote thing WS server
func (c *EriaClient) ConnectThing(url string) (*consumer.ConsumedThing, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal().Err(err).Msg("[eria:ConsumeThing]")
		return nil, errors.New("can't open thing url")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal().Str("status", resp.Status).Str("url", url).Msg("[eria:ConsumeThing] incorrect response")
		return nil, errors.New("incorrect response")
	}

	var td thing.Thing
	if err := json.NewDecoder(resp.Body).Decode(&td); err != nil {
		log.Fatal().Str("url", url).Err(err).Msg("[eria:ConsumeThing]")
		return nil, errors.New("can't decode json")
	}

	consumedThing := c.Consume(&td)
	c.things[td.ID] = consumedThing
	if err != nil {
		return nil, err
	}
	return consumedThing, nil
}
