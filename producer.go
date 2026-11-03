package eria

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/project-eria/eria-core/model"
	eriaproducer "github.com/project-eria/eria-core/producer"
	"github.com/project-eria/go-wot/auth"
	"github.com/project-eria/go-wot/guard"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
	zlog "github.com/rs/zerolog/log"
)

var (
	eriaProducer *eriaproducer.EriaProducer
)

func Producer(instance string) *eriaproducer.EriaProducer {
	if eriaProducer == nil {
		zlog.Trace().Msg("[core:GetProducer] Creating producer")
		eriaProducer = eriaproducer.New(eriaConfig.Host, eriaConfig.Port, eriaConfig.ExposedAddr, _appName, AppVersion, BuildDate)
		eriaProducer.SetMdns(!eriaConfig.Mdns.Disabled, eriaConfig.Mdns.Scheme, eriaConfig.Network)
		cfg, _ := guardConfig()
		g, err := guard.New(cfg)
		if err != nil {
			zlog.Fatal().Err(err).Msg("[core:Producer] Can't load the auth keyset")
		}
		eriaProducer.SetAuthGuard(g)
		AddAppController(eriaProducer, instance)
	}
	return eriaProducer
}

// DeployPath is set by an app whose Things check the master's signature on
// every command themselves (the agent): its guard asks no token.
var DeployPath bool

// guardConfig is the app's guard, and the master's URL for the inventory:
// the `guard:` block, or, without a keyset there, the settings the agent
// wrote beside the app key on the master's order (AUTH-DESIGN §4.7 step 5c).
func guardConfig() (guard.Config, string) {
	a := eriaConfig.Guard
	cfg := guard.Config{Channel: a.Channel, Deployment: a.Deployment, Keyset: a.Keyset, Roots: a.Roots,
		AppKey: appKey(), Floor: timeFloor(), DeployPath: DeployPath}
	master := a.Master
	if path := os.Getenv(auth.AppKeyEnv); a.Keyset == "" && path != "" {
		s, err := guard.ReadSettings(path)
		switch {
		case err == nil:
			cfg.Settings, cfg.Channel, master = &s, s.Channel, s.Master
		case !errors.Is(err, os.ErrNotExist):
			zlog.Fatal().Err(err).Msg("[core:guardConfig] Can't read the activation written by the agent")
		}
	}
	return cfg, master
}

// timeFloor keeps the signed time's floor (AUTH-DESIGN §4.3) beside the app
// key, in the app's data folder; "" without a key (an app the agent did
// not install): no floor kept.
func timeFloor() string {
	if path := os.Getenv(auth.AppKeyEnv); path != "" {
		return filepath.Join(filepath.Dir(path), "time-floor")
	}
	return ""
}

func NewThingDescription(urn string, tdVersion string, title string, description string, capabilities []string) (*thing.Thing, error) {
	td, err := model.NewThingFromModels(
		urn,
		tdVersion,
		title,
		description,
		capabilities,
	)

	if err != nil {
		return nil, err
	}

	td.AddContext("schema", "https://schema.org/")

	// Add Security
	noSecurityScheme := securityScheme.NewNoSecurity()
	td.AddSecurity("no_sec", noSecurityScheme)

	return td, err
}
