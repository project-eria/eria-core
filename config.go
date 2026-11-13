package eria

import (
	"github.com/project-eria/eria-core/automations"
	configmanager "github.com/project-eria/eria-core/config-manager"
	zlog "github.com/rs/zerolog/log"
)

var eriaConfig = struct {
	// Network names the ERIA network the app belongs to ("" = the default
	// network); ERIA_NETWORK overrides it, see resolveNetwork
	Network     string `yaml:"network,omitempty"`
	Host        string `yaml:"host,omitempty"`
	Port        uint   `yaml:"port" default:"80"`
	Location    string `yaml:"location" required:"true"`
	ExposedAddr string `yaml:"exposedAddr,omitempty"`
	// TLS serves the Things over HTTPS with this certificate, signed by the
	// deployment key (the coordinator's: eria-coordinator tls-self, AUTH-DESIGN §4.7
	// step 5d); announced as https.
	TLS struct {
		Cert string `yaml:"cert,omitempty"`
		Key  string `yaml:"key,omitempty"`
	} `yaml:"tls,omitempty"`
	Automations []automations.AutomationConfig `yaml:"automations,omitempty"`
	ContextsRef string                         `yaml:"contextsRef,omitempty"`
	Mdns        struct {
		Disabled         bool   `yaml:"disabled,omitempty"`
		Scheme           string `yaml:"scheme,omitempty" default:"http"`
		DiscoveryTimeout string `yaml:"discoveryTimeout,omitempty" default:"60s"`
	} `yaml:"mdns,omitempty"`
	RemoteThings []RemoteThingConfig `yaml:"remoteThings,omitempty"`
	// Guard checks the eria-auth tokens on the exposed Things (AUTH-DESIGN
	// §2.7): without a keyset file, generation 0 (nosec, with a warning).
	Guard struct {
		Channel     string `yaml:"channel,omitempty"`     // eria-auth Things channel: http://auth:8444
		Deployment  string `yaml:"deployment,omitempty"`  // aud of the tokens
		Keyset      string `yaml:"keyset,omitempty"`      // auth keyset signed by the deployment key
		Roots       string `yaml:"roots,omitempty"`       // deployment roots, PEM
		Coordinator string `yaml:"coordinator,omitempty"` // coordinator's TD URL: its inventory says where the service token may go (F037)
	} `yaml:"guard,omitempty"`
}{}

// RemoteThingConfig declares a remote Thing to consume, either by direct URL
// of its Thing Description, or by DNS-SD identifier (discovered on
// the app's network). Exactly one of Url/Mdns must be set.
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
