package eria

import (
	"time"

	"github.com/project-eria/eria-core/automations"
	zlog "github.com/rs/zerolog/log"
)

func StartAutomations(instance string) {
	if eriaConfig.Automations != nil {
		location, err := time.LoadLocation(eriaConfig.Location)
		if err != nil {
			zlog.Error().Err(err).Msg("[core:StartAutomation]")
			return
		}
		exposedThings := Producer(instance).GetThings()
		automations.Start(location, eriaConfig.Automations, eriaConfig.ContextsRef, exposedThings, _consumedThings)
	} else {
		zlog.Info().Msg("[core:StartAutomation] No automations found, skipping...")
	}
}
