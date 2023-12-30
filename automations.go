package eria

import (
	"github.com/project-eria/eria-core/automations"
	zlog "github.com/rs/zerolog/log"
)

func startAutomations(instance string) {
	if eriaConfig.Automations != nil {
		exposedThings := Producer(instance).GetThings()
		cronScheduler := GetCronScheduler()
		automations.Start(_location, eriaConfig.Automations, eriaConfig.ContextsRef, exposedThings, _consumedThings, cronScheduler)
	} else {
		zlog.Info().Msg("[core:StartAutomation] No automations found, skipping...")
	}
}
