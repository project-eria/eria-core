package eria

import (
	"github.com/go-co-op/gocron/v2"
	zlog "github.com/rs/zerolog/log"
)

var (
	_cronScheduler gocron.Scheduler
)

func GetCronScheduler() gocron.Scheduler {
	if _cronScheduler == nil {
		scheduler, err := gocron.NewScheduler(
			gocron.WithLocation(_location),
		)
		if err != nil {
			zlog.Fatal().Err(err).Msg("[core:GetCronScheduler] Failed to create cron scheduler")
		}
		_cronScheduler = scheduler
	}

	return _cronScheduler
}

func startCronScheduler() {
	if _cronScheduler != nil {
		_cronScheduler.Start()
	}
}

func stopCronScheduler() {
	if _cronScheduler != nil {
		_cronScheduler.Shutdown()
	}
}
