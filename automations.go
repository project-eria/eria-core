package eria

import (
	"sync"
	"time"

	"github.com/project-eria/eria-core/automations"
	eriaproducer "github.com/project-eria/eria-core/producer"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

func StartAutomations(eriaProducer *eriaproducer.EriaProducer) {
	if eriaConfig.Automations != nil {
		location, err := time.LoadLocation(eriaConfig.Location)
		if err != nil {
			zlog.Error().Err(err).Msg("[core:StartAutomation]")
			return
		}
		now := time.Now().In(location)

		// Connects contexts thing (Should be moved to general connections)
		zlog.Info().Str("url", eriaConfig.ContextsUrl).Msg("[core:StartAutomation] Connecting remote contexts Thing")
		eriaConsumer := GetConsumer()
		var (
			contextsThing consumer.ConsumedThing
			wg            = &sync.WaitGroup{}
		)
		wg.Add(1)
		eriaConsumer.ConnectThing(eriaConfig.ContextsUrl, func(t consumer.ConsumedThing) {
			contextsThing = t
			wg.Done()
		}, func(err error) {
			zlog.Fatal().Err(err).Msg("[main] Can't connect context Thing")
		})

		wg.Wait()
		exposedThings := eriaProducer.GetThings()
		automations.Start(now, eriaConfig.Automations, contextsThing, exposedThings)
	} else {
		zlog.Info().Msg("[core:StartAutomation] No automations found, skipping...")
	}
}
