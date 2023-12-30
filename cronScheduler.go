package eria

import "github.com/go-co-op/gocron"

var (
	_cronScheduler *gocron.Scheduler
)

func GetCronScheduler() *gocron.Scheduler {
	if _cronScheduler == nil {
		_cronScheduler = gocron.NewScheduler(_location)
	}

	return _cronScheduler
}

func startCronScheduler() {
	if _cronScheduler != nil {
		_cronScheduler.StartAsync()
	}
}
