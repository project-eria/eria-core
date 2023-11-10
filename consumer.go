package eria

import (
	eriaconsumer "github.com/project-eria/eria-core/consumer"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

var (
	eriaConsumer *eriaconsumer.EriaConsumer
)

func GetConsumer() *eriaconsumer.EriaConsumer {
	if eriaConsumer == nil {
		zlog.Trace().Msg("[core:GetConsumer] Creating consumer")
		eriaConsumer = eriaconsumer.New()
	}
	return eriaConsumer
}

func ConnectThing(url string, onConnected func(*consumer.ConsumedThing), onError func(error)) {
	eriaConsumer := GetConsumer()
	eriaConsumer.ConnectThing(url, onConnected, onError)
}
