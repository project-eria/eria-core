package eria

import (
	"github.com/project-eria/eria-core/automations"
	zlog "github.com/rs/zerolog/log"
)

func startAutomations(instance string) {
	if eriaConfig.Automations != nil {
		exposedThings := Producer(instance).GetThings()
		cronScheduler := GetCronScheduler()
		consumer := Consumer()
		automations.Start(_location, eriaConfig.Automations, exposedThings, consumer, cronScheduler)
	} else {
		zlog.Info().Msg("[core:StartAutomation] No automations found, skipping...")
	}
}
