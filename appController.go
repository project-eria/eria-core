package eria

import (
	"errors"

	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func (s *EriaServer) AddAppController(instance string) {
	eriaAppControllerTd, err := thing.New("eria:app:controler:"+instance, CoreVersion, "EriaAppController", "Eria App Controller", nil)
	if err != nil {
		zlog.Panic().Err(err).Msg("[core] Can't create app controller thing")
	}

	if err := AddModel(eriaAppControllerTd, "EriaAppController", ""); err != nil {
		zlog.Panic().Err(err).Msg("[core] Can't add controller thing model")
	}

	// Add Security
	noSecurityScheme := securityScheme.NewNoSecurity()
	eriaAppControllerTd.AddSecurity("nosec_sc", noSecurityScheme)

	eriaAppControllerThing, err := s.AddThing("eria", eriaAppControllerTd)
	if err != nil {
		zlog.Panic().Err(err).Msg("[core] Can't add app controller")
	}

	eriaAppControllerThing.SetPropertyReadHandler("logLevel", logLevelRead)
	eriaAppControllerThing.SetPropertyWriteHandler("logLevel", logLevelWrite)
}

func logLevelRead(t *producer.ExposedThing, name string, options map[string]string) (interface{}, error) {
	return _logLevel.String(), nil
}

func logLevelWrite(t *producer.ExposedThing, name string, value interface{}, options map[string]string) error {
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
