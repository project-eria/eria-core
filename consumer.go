package eria

import (
	"context"
	"sync"
	"time"

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

// ConnectThings connects the remote Things declared in the config
// (`remoteThings`). Things declared `onConnect: required` (default) block
// until connected — the app won't start without them; `background` ones
// connect whenever they show up.
func ConnectThings() {
	validateRemoteThings()
	discoveryTimeout := mdnsDiscoveryTimeout()
	wg := &sync.WaitGroup{}
	for _, remoteThing := range eriaConfig.RemoteThings {
		remoteThing := remoteThing // Copy https://go.dev/doc/faq#closures_and_goroutines
		required := remoteThing.OnConnect != "background"
		if required {
			wg.Add(1)
			go connectRemoteThing(remoteThing, wg, discoveryTimeout)
		} else {
			go connectRemoteThing(remoteThing, nil, 0)
		}
	}
	wg.Wait()
}

// connectRemoteThing resolves (mDNS) then connects one remote Thing, applying
// its onConnect/onLoss policies. wg is non-nil for `onConnect: required`
// Things; a nil wg means background mode (never fatal, no deadline).
func connectRemoteThing(remoteThing RemoteThingConfig, wg *sync.WaitGroup, discoveryTimeout time.Duration) {
	url := remoteThing.Url
	if remoteThing.Mdns != "" {
		discovered, err := eriaconsumer.DiscoverMdns(context.Background(), remoteThing.Mdns, discoveryTimeout)
		if err != nil {
			// Only reachable with a deadline, i.e. in required mode
			zlog.Fatal().Err(err).Str("mdns", remoteThing.Mdns).Msg("[core:ConnectThings] Can't discover remote Thing")
		}
		url = discovered
	}

	var connected sync.Once
	opts := []eriaconsumer.ThingOption{
		eriaconsumer.WithThingTags(remoteThing.Tags),
		eriaconsumer.WithThingMdnsInstance(remoteThing.Mdns),
		eriaconsumer.WithThingOnConnected(func(t *eriaconsumer.EriaThing) {
			// Also called after each reconnection: release the WaitGroup once
			if wg != nil {
				connected.Do(wg.Done)
			}
		}),
		eriaconsumer.WithThingOnError(func(t *eriaconsumer.EriaThing, err error) {
			if wg != nil {
				zlog.Fatal().Err(err).Msg("[core:ConnectThings] Can't connect remote Thing")
			}
			zlog.Warn().Err(err).Str("url", url).Msg("[core:ConnectThings] Can't connect remote Thing (background), retrying")
			go retryBackground(remoteThing)
		}),
		eriaconsumer.WithThingOnDisconnected(onLossHandler(remoteThing)),
	}
	ConnectThing(url, opts...)
}

// retryBackground retries a background connection after a pause (the
// discovery/connect cycle restarts from scratch, so a moved Thing is found).
func retryBackground(remoteThing RemoteThingConfig) {
	time.Sleep(30 * time.Second)
	connectRemoteThing(remoteThing, nil, 0)
}

// onLossHandler returns the callback applying the `onLoss` policy when an
// interaction with the Thing fails.
func onLossHandler(remoteThing RemoteThingConfig) func(*eriaconsumer.EriaThing, error) {
	return func(t *eriaconsumer.EriaThing, err error) {
		switch remoteThing.OnLoss {
		case "fatal":
			zlog.Fatal().Err(err).Str("mdns", remoteThing.Mdns).Str("url", remoteThing.Url).Msg("[core:ConnectThings] Remote Thing lost")
		case "ignore":
			zlog.Warn().Err(err).Str("mdns", remoteThing.Mdns).Str("url", remoteThing.Url).Msg("[core:ConnectThings] Remote Thing lost (ignored)")
		default: // rediscover
			zlog.Warn().Err(err).Str("mdns", remoteThing.Mdns).Str("url", remoteThing.Url).Msg("[core:ConnectThings] Remote Thing lost, rediscovering")
			url := remoteThing.Url
			if remoteThing.Mdns != "" {
				discovered, err := eriaconsumer.DiscoverMdns(context.Background(), remoteThing.Mdns, 0)
				if err != nil {
					zlog.Error().Err(err).Str("mdns", remoteThing.Mdns).Msg("[core:ConnectThings] rediscovery failed")
					return
				}
				url = discovered
			}
			Consumer().ReconnectThing(t, url)
		}
	}
}

// validateRemoteThings enforces what config-manager tags can't express:
// url XOR mdns, and the onConnect/onLoss enums.
func validateRemoteThings() {
	for i, remoteThing := range eriaConfig.RemoteThings {
		if (remoteThing.Url == "") == (remoteThing.Mdns == "") {
			zlog.Fatal().Int("index", i).Msg("[core:ConnectThings] remoteThings entry must set exactly one of url/mdns")
		}
		switch remoteThing.OnConnect {
		case "", "required", "background":
		default:
			zlog.Fatal().Str("onConnect", remoteThing.OnConnect).Msg("[core:ConnectThings] onConnect must be required|background")
		}
		switch remoteThing.OnLoss {
		case "", "rediscover", "ignore", "fatal":
		default:
			zlog.Fatal().Str("onLoss", remoteThing.OnLoss).Msg("[core:ConnectThings] onLoss must be rediscover|ignore|fatal")
		}
	}
}

func mdnsDiscoveryTimeout() time.Duration {
	timeout, err := time.ParseDuration(eriaConfig.Mdns.DiscoveryTimeout)
	if err != nil || timeout <= 0 {
		zlog.Warn().Str("discoveryTimeout", eriaConfig.Mdns.DiscoveryTimeout).Msg("[core:ConnectThings] invalid mdns.discoveryTimeout, using 60s")
		return 60 * time.Second
	}
	return timeout
}
