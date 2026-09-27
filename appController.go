package eria

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/project-eria/eria-core/automations"
	"github.com/project-eria/eria-core/model"
	propertyModel "github.com/project-eria/eria-core/model/property"
	eriaproducer "github.com/project-eria/eria-core/producer"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

// AddAppController expose, sur le producer de l'application, un Thing WoT
// « EriaAppController » permettant de piloter l'application à distance.
//
// Le Thing est enregistré avec la référence "eria", il est donc servi sous le
// préfixe d'URL « /eria ». Les endpoints suivants sont exposés (host:port =
// ceux configurés via `host`/`port`) :
//
//	Thing Description
//	  GET  /eria                 → Thing Description JSON du contrôleur
//
//	Propriété « logLevel » (niveau de log zerolog, observable)
//	  GET  /eria/logLevel        → lit le niveau de log courant
//	                               réponse: "error" | "warn" | "info" | "debug" | "trace"
//	  PUT  /eria/logLevel        → change le niveau de log
//	                               corps JSON: "error" | "warn" | "info" | "debug" | "trace"
//	  WS   /eria/logLevel        → observe les changements du niveau de log
//	                               (observeproperty / unobserveproperty)
//
//	Propriété « automations » (lecture seule)
//	  GET  /eria/automations     → liste JSON des automations planifiées
//	                               [{automation, devices, action, value, schedule, nextRun}]
//
// Les propriétés sont déclarées sans sécurité (`nosec_sc`).
func AddAppController(p *eriaproducer.EriaProducer, instance string) {
	eriaAppControllerTd, err := thing.New(controllerID(instance), CoreVersion, "EriaAppController", "Eria App Controller", nil)
	if err != nil {
		zlog.Panic().Err(err).Msg("[core] Can't create app controller thing")
	}

	if err := model.AddModel(eriaAppControllerTd, "EriaAppController", ""); err != nil {
		zlog.Panic().Err(err).Msg("[core] Can't add controller thing model")
	}

	// Read-only property exposing the list of scheduled automations
	if _, err := model.AddProperty(eriaAppControllerTd, "automations", propertyModel.Automations, nil); err != nil {
		zlog.Panic().Err(err).Msg("[core] Can't add automations property")
	}

	// Add Security
	noSecurityScheme := securityScheme.NewNoSecurity()
	eriaAppControllerTd.AddSecurity("nosec_sc", noSecurityScheme)

	eriaAppControllerThing, err := p.AddThing("eria", eriaAppControllerTd)
	if err != nil {
		zlog.Panic().Err(err).Msg("[core] Can't add app controller")
	}

	eriaAppControllerThing.SetPropertyReadHandler("logLevel", logLevelRead)
	eriaAppControllerThing.SetPropertyWriteHandler("logLevel", logLevelWrite)
	eriaAppControllerThing.SetPropertyReadHandler("automations", automationsRead)
}

// controllerID builds the controller Thing id, also its mDNS instance name.
// It must be unique on the network: `instance` is empty for most apps and
// the same app runs on several hosts (one agent per node), so the id adds
// "<binary>@<host>". The "eria:app:controler:" prefix (sic) is kept: the
// Home Assistant integration recognizes controllers by it.
func controllerID(instance string) string {
	host, _ := os.Hostname()
	host, _, _ = strings.Cut(host, ".") // short name: dots are escaped in DNS-SD
	bin := strings.ReplaceAll(filepath.Base(os.Args[0]), ".", "-")
	id := "eria:app:controler:" + bin + "@" + host
	if instance != "" {
		id += ":" + instance
	}
	return id
}

func automationsRead(t producer.ExposedThing, name string, parameters map[string]interface{}) (interface{}, error) {
	return automations.Scheduled(), nil
}

func logLevelRead(t producer.ExposedThing, name string, parameters map[string]interface{}) (interface{}, error) {
	return _logLevel.String(), nil
}

func logLevelWrite(t producer.ExposedThing, name string, value interface{}, parameters map[string]interface{}) error {
	logLevelStr := value.(string)
	logLevel, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		return errors.New("can't parse log level")
	}
	_logLevel = logLevel
	zerolog.SetGlobalLevel(logLevel)
	zlog.Info().Stringer("log level", logLevel).Msg("[core:AppController] Set log level")
	return nil
}
