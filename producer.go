package eria

import (
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
		a := eriaConfig.Guard
		g, err := guard.New(guard.Config{Channel: a.Channel, Deployment: a.Deployment, Keyset: a.Keyset, Roots: a.Roots, AppKey: appKey(), Floor: timeFloor()})
		if err != nil {
			zlog.Fatal().Err(err).Msg("[core:Producer] Can't load the auth keyset")
		}
		eriaProducer.SetAuthGuard(g)
		AddAppController(eriaProducer, instance)
	}
	return eriaProducer
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
