package eria

import (
	"errors"

	"github.com/project-eria/eria-core/model"
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
// La propriété logLevel est déclarée sans sécurité (`nosec_sc`).
func AddAppController(p *eriaproducer.EriaProducer, instance string) {
	eriaAppControllerTd, err := thing.New("eria:app:controler:"+instance, CoreVersion, "EriaAppController", "Eria App Controller", nil)
	if err != nil {
		zlog.Panic().Err(err).Msg("[core] Can't create app controller thing")
	}

	if err := model.AddModel(eriaAppControllerTd, "EriaAppController", ""); err != nil {
		zlog.Panic().Err(err).Msg("[core] Can't add controller thing model")
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
