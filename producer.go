package eria

import (
	"github.com/project-eria/eria-core/model"
	eriaproducer "github.com/project-eria/eria-core/producer"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
	zlog "github.com/rs/zerolog/log"
)

var (
	eriaProducer *eriaproducer.EriaProducer
)

func GetProducer(instance string) *eriaproducer.EriaProducer {
	if eriaProducer == nil {
		zlog.Trace().Msg("[core:GetProducer] Creating producer")
		eriaProducer = eriaproducer.New(eriaConfig.Host, eriaConfig.Port, eriaConfig.ExposedAddr, _appName, AppVersion, BuildDate)
		AddAppController(eriaProducer, instance)
	}
	return eriaProducer
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
