package eria

import (
	"sync"

	eriaconsumer "github.com/project-eria/eria-core/consumer"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

var (
	eriaConsumer *eriaconsumer.EriaConsumer
)

func Consumer() *eriaconsumer.EriaConsumer {
	if eriaConsumer == nil {
		zlog.Trace().Msg("[core:GetConsumer] Creating consumer")
		eriaConsumer = eriaconsumer.New()
	}
	return eriaConsumer
}

func ConnectThing(url string, onConnected func(consumer.ConsumedThing), onError func(error)) {
	Consumer().ConnectThing(url, onConnected, onError)
}

func ConnectThings() {
	wg := &sync.WaitGroup{}
	_consumedThings = make(map[string]consumer.ConsumedThing)
	for ref, thingUrl := range eriaConfig.RemoteThings {
		wg.Add(1)
		Consumer().ConnectThing(thingUrl, func(t consumer.ConsumedThing) {
			_consumedThings[ref] = t
			wg.Done()
		}, func(err error) {
			zlog.Fatal().Err(err).Msg("[core:ConnectRemoteThings] Can't connect remote Thing")
		})
	}
	wg.Wait()
}
