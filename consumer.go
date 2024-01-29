package eria

import (
	"sync"

	eriaconsumer "github.com/project-eria/eria-core/consumer"
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

func ConnectThing(url string, opts ...eriaconsumer.ThingOption) {
	Consumer().ConnectThing(url, opts...)
}

func ConnectThings() {
	wg := &sync.WaitGroup{}
	wg.Add(len(eriaConfig.RemoteThings))
	for _, remoteThing := range eriaConfig.RemoteThings {
		remoteThing := remoteThing // Copy https://go.dev/doc/faq#closures_and_goroutines
		ConnectThing(remoteThing.Url,
			eriaconsumer.WithThingTags(remoteThing.Tags),
			eriaconsumer.WithThingOnConnected(func(t *eriaconsumer.EriaThing) {
				wg.Done()
			}),
			eriaconsumer.WithThingOnError(func(t *eriaconsumer.EriaThing, err error) {
				zlog.Fatal().Err(err).Msg("[core:ConnectRemoteThings] Can't connect remote Thing")
				wg.Done()
			}))
	}
	wg.Wait()
}
