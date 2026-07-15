package eria

import (
	"github.com/project-eria/eria-core/automations"
	configmanager "github.com/project-eria/eria-core/config-manager"
	zlog "github.com/rs/zerolog/log"
)

var eriaConfig = struct {
	Host         string                         `yaml:"host,omitempty"`
	Port         uint                           `yaml:"port" default:"80"`
	Location     string                         `yaml:"location" required:"true"`
	ExposedAddr  string                         `yaml:"exposedAddr,omitempty"`
	Automations  []automations.AutomationConfig `yaml:"automations,omitempty"`
	ContextsRef  string                         `yaml:"contextsRef,omitempty"`
	Mdns         struct {
		Disabled         bool   `yaml:"disabled,omitempty"`
		Scheme           string `yaml:"scheme,omitempty" default:"http"`
		DiscoveryTimeout string `yaml:"discoveryTimeout,omitempty" default:"60s"`
	} `yaml:"mdns,omitempty"`
	RemoteThings []RemoteThingConfig `yaml:"remoteThings,omitempty"`
}{}

// RemoteThingConfig declares a remote Thing to consume, either by direct URL
// of its Thing Description, or by DNS-SD identifier (discovered on
// `_wot._tcp`). Exactly one of Url/Mdns must be set.
type RemoteThingConfig struct {
	Url  string   `yaml:"url,omitempty"`
	Mdns string   `yaml:"mdns,omitempty"` // DNS-SD instance (TD id without "urn:") or mDNS host label
	Tags []string `yaml:"tags,omitempty"`
	// OnConnect: "required" (default) blocks the app start and is fatal if
	// the Thing can't be found/connected; "background" lets the app start
	// and connects whenever the Thing shows up.
	OnConnect string `yaml:"onConnect,omitempty" default:"required"`
	// OnLoss: "rediscover" (default) re-resolves (mDNS) or re-fetches (URL)
	// the Thing when an interaction fails; "ignore" logs and stays
	// disconnected; "fatal" stops the app.
	OnLoss string `yaml:"onLoss,omitempty" default:"rediscover"`
}

// loadConfig Loads the config file into a struct
func loadConfig(config interface{}) *configmanager.ConfigManager {
	cm, err := configmanager.Init(*_configPath, config, &eriaConfig)
	if err != nil {
		if configmanager.IsFileMissing(err) {
			zlog.Fatal().Str("filePath", *_configPath).Msg("[core:LoadConfig] Config file do not exists...")
		} else {
			zlog.Fatal().Str("filePath", *_configPath).Err(err).Msg("[core:LoadConfig]")
		}
	}

	if err := cm.Load(); err != nil {
		zlog.Fatal().Err(err).Msg("[core:LoadConfig]")
	}
	return cm
}
