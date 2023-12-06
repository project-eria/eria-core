package eria

import (
	"time"

	"github.com/project-eria/eria-core/automations"
	zlog "github.com/rs/zerolog/log"
)

func StartAutomations() {
	if eriaConfig.Automations != nil {
		location, err := time.LoadLocation(eriaConfig.Location)
		if err != nil {
			zlog.Error().Err(err).Msg("[core:StartAutomation]")
			return
		}
		eriaProducer := GetProducer("")
		exposedThings := eriaProducer.GetThings()
		automations.Start(location, eriaConfig.Automations, eriaConfig.ContextsRef, exposedThings, _consumedThings)
	} else {
		zlog.Info().Msg("[core:StartAutomation] No automations found, skipping...")
	}
}
